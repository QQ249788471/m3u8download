package task

import (
	"fmt"
	"github.com/lyd2/goff-net-task/bean"
	"github.com/lyd2/goff-net-task/util"
	"github.com/lyd2/goff-net-task/worker/m3u8ex"
	"github.com/sirupsen/logrus"
)

var tsTaskExit chan struct{}

type DownloadMode int

const (
	// 首次下载
	TS_DOWNLOAD_FIRST DownloadMode = 1

	// 重试下载
	TS_DOWNLOAD_RETRY DownloadMode = 2
)

func init() {

	tsTaskExit = make(chan struct{})

	// 开始运行下载 ts 的 task
	// 它会一直运行，而不会退出，除非调用 TsTaskExit
	for i := 0; i < util.Conf.TsDownloadCount; i++ {
		downloadTsTask(TS_DOWNLOAD_FIRST, fmt.Sprintf("f_%d", i))
	}

	// 重试下载 ts
	for i := 0; i < util.Conf.TsRetryDownloadCount; i++ {
		downloadTsTask(TS_DOWNLOAD_RETRY, fmt.Sprintf("r_%d", i))
	}

}

// 退出 ts 下载任务
func TsTaskExit() {
	close(tsTaskExit)
}

// ts 下载任务
func downloadTsTask(mode DownloadMode, gid string) {

	go func() {

		for {

			select {
			case <-tsTaskExit:
				// 退出协程
				return
			default:
				// 继续处理任务
			}

			// 查找下载任务
			var recordInfo *bean.RecordInfo
			if mode == TS_DOWNLOAD_FIRST {
				// 首次下载模式
				// logrus.Infof("[Thread %s] CompareAndSwap_First START", gid)
				recordInfo = (bean.RecordInfo{}).CompareAndSwap_First(gid)
				if recordInfo == nil {
					// logrus.Infof("[Thread %s] CompareAndSwap_First FALSE", gid)
					util.Sleep()
					continue
				}
				// logrus.Infof("[Thread %s] CompareAndSwap_First TRUE, TS_ID = %d", gid, recordInfo.Id)
			} else {
				// 重试下载模式
				// logrus.Infof("[Thread %s] CompareAndSwap_Retry START", gid)
				recordInfo = (bean.RecordInfo{}).CompareAndSwap_Retry(gid)
				if recordInfo == nil {
					// logrus.Infof("[Thread %s] CompareAndSwap_Retry FALSE", gid)
					util.Sleep()
					continue
				}
				// logrus.Infof("[Thread %s] CompareAndSwap_Retry TRUE, TS_ID = %d", gid, recordInfo.Id)
			}

			// 开始下载
			size, t, err := m3u8ex.DownloadTs(recordInfo)
			if err != nil {
				// 下载失败
				logrus.Error(err)
				recordInfo.DownloadFail(fmt.Errorf("[ThreadId: %s] %s", gid, err.Error()))
				continue
			}

			// 下载成功
			recordInfo.DownloadSuccess(size, t)
		}

	}()

}
