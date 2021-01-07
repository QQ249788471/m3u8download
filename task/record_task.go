package task

import (
	"github.com/lyd2/goff-net-task/bean"
	"github.com/lyd2/goff-net-task/library/queue"
	"github.com/lyd2/goff-net-task/util"
	"github.com/lyd2/goff-net-task/worker/m3u8ex"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"time"
)

// m3u8 请求间隔时间
const TASK_INTERVAL = 500 * time.Millisecond

// 存储 ts 文件名的队列长度
const QUEUE_LENGTH = 2048

type RecordTask struct {

	// 收录任务
	record bean.Record

	// 收录任务处理节点
	recordNode bean.RecordNode

	// 此任务文件存储路径
	path string

	// 退出状态跟踪
	exitWatcher chan struct{}

	// 存储 ts 文件名的队列
	q *queue.Queue
}

// 新建一个收录任务
func NewRecordTask(record bean.Record, recordNode bean.RecordNode) *RecordTask {
	path, err := getRecordPath(record)
	if err != nil {
		bean.NewLog(recordNode, err.Error()).Error()
		return nil
	}

	return &RecordTask{
		record:     record,
		recordNode: recordNode,
		path:       path,
		q:          queue.New(QUEUE_LENGTH),
	}
}

// 生成收录的文件保存路径
func getRecordPath(record bean.Record) (path string, err error) {

	t := time.Now()

	// 最多重试 10 次
	for i := 0; i < 10; i++ {
		path, err = util.GetTsPath(t, record.Id, true)
		if err == nil {
			break
		}
	}

	return
}

// 开始执行任务
func (t *RecordTask) Start() {

	bean.NewLog(t.recordNode, "TASK START").Info()

	// task loop
	t.m3u8()

	// 执行完成
	t.finish()

}

// 任务已完成
func (t *RecordTask) finish() {
	bean.NewLog(t.recordNode, "TASK END").Info()
	// TODO: 轮询此任务的所有 ts 下载记录，全部完成时才修改状态
	t.recordNode.Finish()
}

// 用户主动退出
func (t *RecordTask) exit() {
	bean.NewLog(t.recordNode, "TASK EXIT").Info()
	t.recordNode.Delete()
}

// 完成 m3u8 的下载解析功能
func (t *RecordTask) m3u8() {

	ticker := time.NewTicker(TASK_INTERVAL)

	for time.Now().Unix() <= int64(t.record.EndTime) {

		// 检测当前任务节点记录的删除状态
		// 在运行时停止这个任务，或者删除这个任务，会触发此条件
		if t.recordNode.IsDelete() {
			// 此节点已经不再需要执行这个任务，可能的原因是程序需要退出，或者任务被停止，或者任务被删除
			logrus.Infof("[RecordTask deleted]: RecordId=%d, RecordNodeId=%d", t.record.Id, t.recordNode.Id)
			return
		}

		select {
		case <-t.exitWatcher:
			// 任务已退出
			// 用户手动关闭或重启程序时，进入此分支
			t.exit()
			return
		case <-ticker.C:
			// 执行任务
			break
		}

		// 下载并解析 m3u8
		m3u8f, err := m3u8ex.DownloadAndParseM3u8(t.record.Live)
		if err != nil {
			// m3u8 下载或解析错误
			bean.NewLog(t.recordNode, m3u8f.String()).Error()
			continue
		}
		// bean.NewLog(t.recordNode, m3u8f.String()).Info()

		// 查询 ts 文件
		tsInfoList := m3u8ex.SearchTs(&m3u8f)

		// 创建 ts row
		for _, v := range tsInfoList {
			v.Filename = filepath.Join(t.path, v.Filename)

			if v.Status {
				// ts uri 生成成功

				if !t.q.Search(v.Filename) {

					// 队列中没有此 ts 的 filename（说明是新的 ts）
					t.q.Insert(v.Filename)
					recordInfo := bean.NewRecordInfo(t.recordNode, v.Uri, v.Filename)
					err := recordInfo.Insert()
					if err != nil {
						// 插入 ts 信息到数据库失败
						bean.NewLog(t.recordNode, err.Error()).Error()
						logrus.Error(err)
						continue
					}
					// fmt.Println(recordInfo)

					// 这里只是生成 ts 文件下载任务（record_info 表中的一行，表示一个 ts 的下载任务）
					// 具体的下载操作，在 ts 的下载任务协程中执行

				}

			} else {
				// ts uri 生成失败
				bean.NewLog(t.recordNode, v.Error).Error()
			}
		}
	}

}
