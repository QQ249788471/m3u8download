package bean

import (
	"github.com/lyd2/goff-net-task/util"
	"time"
)

type TaskStatus int

const (
	DOING TaskStatus = 1
	DONE  TaskStatus = 2
)

type RecordNode struct {
	Id int `gorm:"column:id"`

	// Record.Id
	RecordId int `gorm:"column:record_id"`

	// 收录任务的开始时间
	RecordStartTime int `gorm:"column:record_start_time"`

	// 节点名称
	NodeName string `gorm:"column:node_name"`

	// 节点IP
	NodeIp string `gorm:"column:node_ip"`

	// 任务状态 TaskStatus
	Status TaskStatus `gorm:"column:status"`

	// 删除时间
	DeletedAt int `gorm:"column:deleted_at"`
}

func (RecordNode) TableName() string {
	return "record_node"
}

func NewRecordNode(r *Record) *RecordNode {
	return &RecordNode{
		RecordId:        r.Id,
		RecordStartTime: r.StartTime,
		NodeName:        util.Conf.NodeName,
		NodeIp:          util.ServerIp,
		Status:          DOING,
		DeletedAt:       0,
	}
}

func (rn *RecordNode) Insert() error {
	err := util.Db.Create(rn).Error
	return err
}

func (rn *RecordNode) Get(where interface{}) error {
	return util.Db.Table(rn.TableName()).Where(where).First(rn).Error
}

func (rn *RecordNode) IsDelete() bool {
	recordNode := &RecordNode{}
	util.Db.Table(recordNode.TableName()).
		Where("id = ?", rn.Id).
		Where("deleted_at <> ?", 0).
		Find(recordNode)

	// fmt.Println(recordNode)

	return recordNode.Id > 0
}

// 任务已完成
func (rn *RecordNode) Finish() {
	util.Db.Model(RecordNode{}).
		Where("id = ?", rn.Id).
		Updates(map[string]interface{}{
			"status": DONE,
		})
}

// 删除任务节点记录
// 注意，这个删除状态只表示这个节点不再执行这个任务，但是并不表示这个收录任务已经结束
// 也就是说，如果收录任务在查询相关记录时，即便这个节点记录已被删除，它所下载的 ts 文件也需要被查询出来
func (rn *RecordNode) Delete() {
	util.Db.Model(RecordNode{}).
		Where("id = ?", rn.Id).
		Updates(map[string]interface{}{
			"status":     DONE,
			"deleted_at": time.Now().Unix(),
		})
}
