package util

import (
	"strconv"
	"time"
)

// 获取星期几，1-7 表示周一到周日
func Weekday() string {
	w := time.Now().Weekday()
	if w == time.Weekday(0) {
		return "7"
	}

	return strconv.Itoa(int(w))
}

func Sleep() {
	// 睡眠间隔时间
	time.Sleep(time.Second * time.Duration(Conf.Interval))
}

func TimeFormat(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
