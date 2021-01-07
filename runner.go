package main

import (
	"github.com/lyd2/goff-net-task/bean"
	"github.com/lyd2/goff-net-task/task"
	"github.com/lyd2/goff-net-task/util"
)

func Run() {

	go func() {

		// task loop
		for {

			// 任务轮询间隔时间
			util.Sleep()

			// 查询任务列表
			records := bean.GetRecords()
			if len(records) == 0 {
				// logrus.Info("[RecordTask is empty]")
				continue
			}
			// fmt.Println(records)

			for _, record := range records {

				// 查询收录任务是否已存在
				// (查询未删除的任务节点记录，因为如果某个任务处理节点记录已删除，则此任务可以再被其他节点抢占了)
				// (也就是说，一个任务只能有一个节点正在执行)
				rn := bean.RecordNode{}
				_ = rn.Get(map[string]interface{}{
					"record_id":         record.Id,
					"record_start_time": record.StartTime,
					"deleted_at":        0,
				})
				if rn.Id > 0 {
					// logrus.Infof("[RecordTask existed]: %s", record.String())
					continue
				}

				// 创建任务并写入通道
				recordTask := task.NewRecordTask(record, *bean.NewRecordNode(&record))
				if recordTask != nil {
					task.WriteRecordTask(recordTask)
				}

			}

		}

	}()

}
