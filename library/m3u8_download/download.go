package m3u8download

import (
	"fmt"
	"github.com/lyd2/goff-net-task/util"
	"io"
	"net/http"
	"os"
	"time"
)

// 下载 ts
func DownloadTs(uri, filename string) (size int, t string, err error) {

	// 开始执行时间
	startTime := time.Now()

	// 获取 ts 文件内容
	res, err := util.Client.Get(uri)
	if err != nil {
		return fail(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fail(fmt.Errorf("HTTP %d", res.StatusCode))
	}

	// 创建文件
	fd, err := os.Create(filename)
	if err != nil {
		return fail(err)
	}
	defer fd.Close()

	// 写入 writer
	var n int64
	if n, err = io.Copy(fd, res.Body); err != nil && err != io.ErrUnexpectedEOF {
		return fail(err)
	}

	// 计算时间间隔
	sub := time.Now().Sub(startTime)

	// 计算文件大小
	size = int(n)

	return success(size, sub)
}

// return (size int, t string, err error)
func success(size int, t time.Duration) (int, string, error) {
	return size, fmt.Sprintf("%.2f", float32(t)/float32(time.Second)), nil
}

// return (size int, t string, err error)
func fail(err error) (int, string, error) {
	return 0, "0", err
}
