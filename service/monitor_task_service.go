package service

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"spiderman-agent/common/model"
)

func MonitorTask(s *Service) (err error) {
	rch := s.client.Watch(context.Background(), s.Name, clientv3.WithPrefix())
	for resp := range rch {
		for _, ev := range resp.Events {
			switch ev.Type {
			case clientv3.EventTypePut:
				// monitor task change
				checkTaskMap(ev, s)
			case clientv3.EventTypeDelete:
				//todo someing
			}
		}
	}
	return
}

func checkTaskMap(ev *clientv3.Event, s *Service) {
	remoteInfo := &model.ServiceInfo{}
	err := json.Unmarshal([]byte(ev.Kv.Value), remoteInfo)
	if err != nil {
		return
	}

	for taskId, remoteTask := range remoteInfo.TaskMap {
		has, localTask := taskExist(taskId, s)
		if !has {
			fmt.Println("节点增加任务")
			s.Info.TaskMap[remoteTask.TaskId] = remoteTask
			//todo 运行任务
		}else{
			//check task and deal with task
			taskDetail(remoteTask, localTask)
			//todo 更改任务状态
		}
	}
	return
}

func taskExist(taskId string, s *Service) (bool, *model.TaskInfo) {
	if task, ok := s.Info.TaskMap[taskId]; ok {
		return true, task
	}
	return false, nil
}

func taskDetail(remoteTask, localTask *model.TaskInfo) {
	fmt.Printf("远程任务: %s \n", remoteTask.TaskName)
	fmt.Printf("本地任务: %s \n", localTask.TaskName)
}
