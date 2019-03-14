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
			checkTaskMap(ev, s)
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

	switch ev.Type {
	case clientv3.EventTypePut:
		// monitor task change
		dealWith(remoteInfo, s)
	case clientv3.EventTypeDelete:
		//stop services
		s.Stop()
	}
	return
}

func dealWith(info *model.ServiceInfo, s *Service) {
	for taskId, remoteTask := range info.TaskMap {
		has, localTask := exist(taskId, s)
		if !has {
			fmt.Println("节点增加任务")
			start(remoteTask, s)
		} else {
			fmt.Println("已经存在任务")
			check(remoteTask, localTask, s)
		}
	}
	return
}

func exist(taskId string, s *Service) (bool, *model.TaskInfo) {
	if task, ok := s.Info.TaskMap[taskId]; ok {
		return true, task
	}
	return false, nil
}

func start(remote *model.TaskInfo, s *Service) {
	s.Info.TaskMap[remote.TaskId] = remote
	GetAccountService().init(remote)
}

func check(remote, local *model.TaskInfo, s *Service) {
	if remote.Config.ProgramUpdateAt != local.Config.ProgramUpdateAt {
		switch remote.Config.Status {
		case model.TaskStatus_Stop:
			GetAccountService().stop(local)
		case model.TaskStatus_Start:
			GetAccountService().init(remote)
		}
	}
	s.Info.TaskMap[remote.TaskId] = remote
}
