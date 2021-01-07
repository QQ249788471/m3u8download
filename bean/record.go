package bean

import (
	"encoding/json"
	"github.com/lyd2/goff-net-task/util"
	"time"
)

// 直播类型
type LiveType int

const (
	// 仅一次节目
	ONLY_ONE = 1

	// 周期性节目
	REPEAT = 2

	// 开启状态
	STATUS_ENABLE = 1

	// 关闭状态
	STATUS_DISABLE = 2
)

type Record struct {
	Id int `gorm:"column:id" json:"id"`

	// 直播地址
	Live string `gorm:"column:live" json:"live"`

	// 开始收录时间
	StartTime int `gorm:"column:start_time" json:"start_time"`

	// 结束收录时间
	EndTime int `gorm:"column:end_time" json:"end_time"`

	// 是否为周期性节目 ONLY_ONE | REPEAT
	IsRepeat int `gorm:"column:is_repeat" json:"is_repeat"`

	// 1-7 周一到周日
	Weekday string `gorm:"column:weekday" json:"weekday"`
}

func (Record) TableName() string {
	return "record"
}

// 获取待执行任务
func GetRecords() []Record {

	// 仅一次的列表
	var recordsNr []Record
	// 周期性列表
	var recordsR []Record

	t := time.Now()

	// 查询仅一次节目的列表
	now := t.Unix()
	util.Db.Table(Record{}.TableName()).
		Where("start_time <= ? and end_time > ?", now, now).
		Where("is_repeat = ?", ONLY_ONE).
		Where("status = ?", STATUS_ENABLE).
		Where("deleted_at = ?", 0).
		Find(&recordsNr)

	// 获取周期性节目的列表
	// 周期性节目的日期统一到 2020-01-01
	loc, _ := time.LoadLocation("Local")
	ts := time.Date(2020, time.January, 1, t.Hour(), t.Minute(), t.Second(), 0, loc).Unix()
	util.Db.Table(Record{}.TableName()).
		Where("start_time <= ? and end_time > ?", ts, ts).
		Where("is_repeat = ? and FIND_IN_SET(?, weekday) > 0", REPEAT, util.Weekday()).
		Where("status = ?", STATUS_ENABLE).
		Where("deleted_at = ?", 0).
		Find(&recordsR)

	// 为周期性节目设置正确的开始时间和结束时间
	for k, v := range recordsR {
		snt := time.Unix(int64(v.StartTime), 0)
		ent := time.Unix(int64(v.EndTime), 0)
		recordsR[k].StartTime = int(time.Date(t.Year(), t.Month(), t.Day(), snt.Hour(), snt.Minute(), snt.Second(), 0, loc).Unix())
		recordsR[k].EndTime = int(time.Date(t.Year(), t.Month(), t.Day(), ent.Hour(), ent.Minute(), ent.Second(), 0, loc).Unix())
	}

	return append(recordsNr, recordsR...)
}

func (r *Record) String() string {
	s, _ := json.Marshal(r)
	return string(s)
}
