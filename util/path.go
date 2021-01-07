package util

import (
	"fmt"
	"os"
	"time"
)

// 判断文件夹是否存在
func PathIsExist(path string) bool {

	_, err := os.Stat(path)

	if err != nil {
		if os.IsExist(err) {
			return true
		}

		if os.IsNotExist(err) {
			return false
		}

		return false
	}

	return true

}

func Mkdir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

// 获取日志的存储路径
func GetLogPath(path string) string {
	return fmt.Sprintf("%s/%s", Conf.LogPath, path)
}

// 获取 Ts 的存储路径
func GetTsPath(t time.Time, task int, create bool) (path string, err error) {
	path = fmt.Sprintf("%s/ts/%d/%d/%d/%d", Conf.DownloadPath, t.Year(), t.Month(), t.Day(), task)
	if create && !PathIsExist(path) {
		err = Mkdir(path)
	}
	return
}

// 获取 Flv 的存储路径
func GetFlvPath(t time.Time, task int, create bool) (path string, err error) {
	path = fmt.Sprintf("%s/flv/%d/%d/%d/%d", Conf.DownloadPath, t.Year(), t.Month(), t.Day(), task)
	if create && !PathIsExist(path) {
		err = Mkdir(path)
	}
	return
}

// 创建一个文件，仅写、不存在则创建、追加写，所有用户可读写、不可执行
func OpenFile(filepath string) (file *os.File, err error) {
	return os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
}
