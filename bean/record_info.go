package bean

import (
	"github.com/lyd2/goff-net-task/util"
	"time"
)

type DownloadStatus int

const (
	WAITING            DownloadStatus = 0
	DOWNLOADING        DownloadStatus = 1
	DOWNLOAD_COMPLETED DownloadStatus = 2
	DOWNLOAD_FAILED    DownloadStatus = 3
)

type RecordInfo struct {

	// id
	Id int `gorm:"column:id"`

	// RecordNode.Id
	RecordNodeId int `gorm:"column:record_node_id"`

	// ts Uri
	Uri string `gorm:"column:uri"`

	// 文件保存路径
	Path string `gorm:"column:path"`

	// 任务状态, DOWNLOADING | DOWNLOAD_COMPLETED | DOWNLOAD_FAILED
	Status DownloadStatus `gorm:"column:status"`

	// 重试次数
	Retry int `gorm:"column:retry"`

	// 文件大小，单位字节
	Size int `gorm:"column:size"`

	// 下载耗时，单位秒
	Time string `gorm:"column:time"`
}

func (RecordInfo) TableName() string {
	return "record_info"
}

func NewRecordInfo(rn RecordNode, uri, path string) *RecordInfo {
	return &RecordInfo{
		RecordNodeId: rn.Id,
		Uri:          uri,
		Path:         path,
		Status:       WAITING,
	}
}

func (ri *RecordInfo) Insert() error {
	return util.Db.Create(ri).Error
}

// cas 抢占一个任务（首次下载）
func (RecordInfo) CompareAndSwap_First(gid string) *RecordInfo {

	// id 从 0 开始查找
	id := 0

	// 最多只检索 n 小时之前的任务
	startTime := time.Now().Add(-time.Duration(util.Conf.SearchTime) * time.Hour)

	for i := 0; i < 10; i++ {

		recordInfo := &RecordInfo{}

		// 重试次数为 0，且状态为 WAITING
		util.Db.Table(recordInfo.TableName()).
			Where("id > ?", id).
			Where("created_at > ?", util.TimeFormat(startTime)).
			Where("retry = ?", 0).
			Where("status = ?", WAITING).
			Find(recordInfo)
		if recordInfo.Id <= 0 {
			// 未查询到数据
			return nil
		}

		// test
		// fmt.Printf("Thread_id = %s, ri_id = %d\n", gid, id)

		// test
		// time.Sleep(10 * time.Second)

		// 尝试更新它的状态和重试次数
		res := util.Db.Model(RecordInfo{}).
			Where("id = ?", recordInfo.Id).
			Where("retry = ?", 0).
			Where("status = ?", WAITING).
			Updates(map[string]interface{}{
				"retry":  1,
				"status": DOWNLOADING,
			})

		// test
		// fmt.Println(res.RowsAffected)

		if res.RowsAffected != 1 {
			// 更新失败，开始下一轮更新
			// 新的查询起到要在这条记录之后
			id = recordInfo.Id
			continue
		}

		// 查找成功
		return recordInfo
	}

	// 超过 10 次
	return nil

}

// cas 抢占一个任务（重试下载）
func (RecordInfo) CompareAndSwap_Retry(gid string) *RecordInfo {

	// id 从 0 开始查找
	id := 0

	// 最多只检索 n 小时之前的任务
	startTime := time.Now().Add(-time.Duration(util.Conf.SearchTime) * time.Hour)

	for i := 0; i < 10; i++ {

		recordInfo := &RecordInfo{}

		// 最大重试次数未超过额定数量，且状态为失败
		util.Db.Table(recordInfo.TableName()).
			Where("id > ?", id).
			Where("created_at > ?", util.TimeFormat(startTime)).
			Where("retry > ?", 0).
			Where("retry < ?", util.Conf.MaxRetry).
			Where("status = ?", DOWNLOAD_FAILED).
			Find(recordInfo)
		if recordInfo.Id <= 0 {
			// 未查询到数据
			return nil
		}

		// test
		// fmt.Printf("Thread_id = %s, ri_id = %d\n", gid, id)

		// test
		// time.Sleep(10 * time.Second)

		// 尝试更新它的重试次数
		res := util.Db.Model(RecordInfo{}).
			Where("id = ?", recordInfo.Id).
			Where("retry = ?", recordInfo.Retry).
			Where("status = ?", DOWNLOAD_FAILED).
			Updates(map[string]interface{}{
				"retry":  recordInfo.Retry + 1,
				"status": DOWNLOADING,
			})

		// test
		// fmt.Println(res.RowsAffected)

		if res.RowsAffected != 1 {
			// 更新失败，开始下一轮更新
			// 新的查询起到要在这条记录之后
			id = recordInfo.Id
			continue
		}

		// 查找成功
		return recordInfo
	}

	// 超过 10 次
	return nil
}

/*
	假设 MaxRetry == 3
	节点模式：[status, retry]

						start
					[WAITING, 0]
                          |
                          |  被 First 线程读取
                          |
                    [DOWNLOADING, 1]
                         /  \
                执行成功 /    \ 执行失败
                       /      \
  [DOWNLOAD_COMPLETED, 1]    [DOWNLOAD_FAILED, 1]
                                       |
                                       | 被 Retry 线程读取
                                       |
                             [DOWNLOADING, 2]
                                   /  \
                          执行成功 /    \ 执行失败
                                 /      \
            [DOWNLOAD_COMPLETED, 2]    [DOWNLOAD_FAILED, 2]
                                               |
                                               | 被 Retry 线程读取
                                               |
                                       [DOWNLOADING, 3]
                                            /   \
                                   执行成功 /     \ 执行失败
                                          /       \
                     [DOWNLOAD_COMPLETED, 3]     [DOWNLOAD_FAILED, 3]

*/

func (ri *RecordInfo) DownloadSuccess(size int, t string) {
	util.Db.Model(RecordInfo{}).
		Where("id = ?", ri.Id).
		Updates(map[string]interface{}{
			"status": DOWNLOAD_COMPLETED,
			"size":   size,
			"time":   t,
		})
}

func (ri *RecordInfo) DownloadFail(err error) {

	util.Db.Model(RecordInfo{}).
		Where("id = ?", ri.Id).
		Updates(map[string]interface{}{
			"status": DOWNLOAD_FAILED,
		})

	NewRecordInfoLog(ri, err.Error()).Insert()

}
