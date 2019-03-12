package service

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
)

func MonitorTask(s *Service) (err error) {
	rch := s.client.Watch(context.Background(), s.Name, clientv3.WithPrefix())
	for resp := range rch {
		for _, ev := range resp.Events {
			switch ev.Type {
			case clientv3.EventTypePut:
				info, err := GetServiceInfo(ev)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(info)
				return nil
			case clientv3.EventTypeDelete:

			}
		}
	}
	return
}

func GetServiceInfo(ev *clientv3.Event) (*ServicesInfo, error) {

	fmt.Println(string([]byte(ev.Kv.Value)))
	info := &ServicesInfo{}
	err := json.Unmarshal([]byte(ev.Kv.Value), info)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return info, nil
}
