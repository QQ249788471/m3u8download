package task

import (
	"github.com/lyd2/goff-net-task/util"
	"github.com/sirupsen/logrus"
)

// 在生产者那端，只是预查询任务，它将查询到的任务投递到具体的 task，再由 task 真正的去抢占
// 为什么要这么做呢，原因很简单，假如直接在生产者那端完成抢占，那如果此节点所有的 task 协程都已经被占用，则此新任务会被阻塞
// 而此时假如有其他节点有空余的 task 协程，它们也无法再去抢占了
// 因此使用预查询的方式，让具体的 task 协程去真正的抢占任务，才是最合理的
var recordTask chan *RecordTask

var m3u8TaskExit chan struct{}

func init() {

	recordTask = make(chan *RecordTask)

	m3u8TaskExit = make(chan struct{})

	// 开启 m3u8 task 协程
	for i := 0; i < util.Conf.MaxTaskCount; i++ {
		downloadM3u8Task()
	}

}

// 写入收录任务
func WriteRecordTask(rt *RecordTask) {
	recordTask <- rt
}

// 退出 m3u8 下载任务
func M3u8TaskExit() {
	close(m3u8TaskExit)
}

func downloadM3u8Task() {

	go func() {

		for {

			// 从任务通道读取任务
			rt := <-recordTask

			// 其它节点也可能会查询到此任务，因此需要抢占
			if err := rt.recordNode.Insert(); err != nil {
				// 此任务已被其它节点抢占
				// 极少进入的分支
				/*
						进入此分支的情况如下

						node_1            node_2
					    [task不存在]
					    [CPU时间片结束]
					                      [task不存在]
					                      [创建task]
					                      ...
					                      [CPU时间片结束]
						[创建task]
					    [创建失败..]
				*/
				// 另外，也可能因为任务被阻塞，当再次获取到任务时，大概率会进入此分支
				logrus.Infof("[RecordTask doing]: RecordNode=%d", rt.recordNode.Id)
				continue
			}

			// 让此任务监听协程退出通道
			rt.exitWatcher = m3u8TaskExit

			// 开始执行
			rt.Start()
		}

	}()
}
