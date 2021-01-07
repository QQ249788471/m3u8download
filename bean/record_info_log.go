package bean

import "github.com/lyd2/goff-net-task/util"

type RecordInfoLog struct {

	// id
	Id int `gorm:"column:id"`

	// RecordInfo.Id
	RecordInfoId int `gorm:"column:record_info_id"`

	// 详细内容
	Content string `gorm:"column:content"`

	// 节点名称
	NodeName string `gorm:"column:node_name"`

	// 节点IP
	NodeIp string `gorm:"column:node_ip"`
}

func (RecordInfoLog) TableName() string {
	return "record_info_log"
}

func NewRecordInfoLog(ri *RecordInfo, content string) *RecordInfoLog {
	return &RecordInfoLog{
		RecordInfoId: ri.Id,
		Content:      content,
		NodeName:     util.Conf.NodeName,
		NodeIp:       util.ServerIp,
	}
}

func (ril *RecordInfoLog) Insert() {
	util.Db.Create(ril)
}
