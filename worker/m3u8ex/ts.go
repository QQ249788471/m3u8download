package m3u8ex

import (
	"github.com/lyd2/goff-net-task/bean"
	m3u8download "github.com/lyd2/goff-net-task/library/m3u8_download"
)

// 下载 ts 文件
func DownloadTs(rinfo *bean.RecordInfo) (size int, t string, err error) {
	return m3u8download.DownloadTs(rinfo.Uri, rinfo.Path)
}
