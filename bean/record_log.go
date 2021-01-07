package bean

import "github.com/lyd2/goff-net-task/util"

type LogLevel int

const (
	INFO  LogLevel = 1
	WARN  LogLevel = 2
	ERROR LogLevel = 3
)

type RecordLog struct {
	Id int `gorm:"column:id"`

	// RecordNode.Id
	RecordNodeId int `gorm:"column:record_node_id"`

	// 日志级别 LogLevel
	Level LogLevel `gorm:"column:level"`

	// 内容
	Content string `gorm:"column:content"`
}

func (RecordLog) TableName() string {
	return "record_log"
}

func NewLog(rn RecordNode, content string) *RecordLog {
	return &RecordLog{
		RecordNodeId: rn.Id,
		Content:      content,
	}
}

func (rl *RecordLog) Info() {
	rl.Level = INFO
	util.Db.Create(rl)
}

func (rl *RecordLog) Warn() {
	rl.Level = WARN
	util.Db.Create(rl)
}

func (rl *RecordLog) Error() {
	rl.Level = ERROR
	util.Db.Create(rl)
}
