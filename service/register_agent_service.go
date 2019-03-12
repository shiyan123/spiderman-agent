package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"log"
	"spiderman-agent/app"
	"spiderman-agent/common/health"
	"spiderman-agent/common/model"
	"spiderman-agent/utils"
	"time"
)

type ServicesInfo struct {
	IP          string
	Name        string
	MachineInfo *health.MachineHealth
	TaskMap     map[string]*model.TaskInfo
}

type Service struct {
	Name    string
	Info    *ServicesInfo
	stop    chan error
	leaseId clientv3.LeaseID
	client  *clientv3.Client
}

func RegisterService() (*Service, error) {
	// create etcd client
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   app.GetApp().Config.CenterServer.Ips,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	s := &Service{
		Name:   fmt.Sprintf("services/%s/%s", app.GetApp().Config.Server.RegisterName, utils.EncodeTimeStr(time.Now())),
		Info:   initServicesInfo(),
		stop:   make(chan error),
		client: cli,
	}

	// register service to etcd
	leaseId, err := s.register()
	if err != nil {
		return nil, err
	}

	//keepAlive service
	go func() {
		s.StartKeepAlive(leaseId)
	}()
	return s, nil
}

func (s *Service) register() (leaseId clientv3.LeaseID, err error) {
	machineInfo, err := health.GetMachineHealth()
	if err != nil {
		return
	}
	s.Info.MachineInfo = machineInfo

	info := &s.Info
	value, err := json.Marshal(info)
	if err != nil {
		return
	}
	resp, err := s.client.Grant(context.TODO(), 5)
	if err != nil {
		return
	}
	_, err = s.client.Put(context.TODO(), s.Name, string(value), clientv3.WithLease(resp.ID))
	if err != nil {
		return
	}
	return resp.ID, nil
}

func initServicesInfo() *ServicesInfo {
	return &ServicesInfo{
		IP:   app.GetApp().Config.Server.RegisterIp,
		Name: app.GetApp().Config.Server.RegisterName,
	}
}

func (s *Service) StartKeepAlive(leaseId clientv3.LeaseID) error {

	ch, err := s.client.KeepAlive(context.TODO(), leaseId)
	if err != nil {
		return err
	}

	//todo 待优化
	for {
		select {
		case err := <-s.stop:
			s.revoke()
			return err
		case <-s.client.Ctx().Done():
			return errors.New("server closed")
		case _, ok := <-ch:
			if !ok {
				log.Println("keep alive channel closed")
				s.revoke()
				return nil
			} else {
				body, _ := json.Marshal(s.Info)
				fmt.Println(string(body))
			}
		}
	}
}

func (s *Service) Stop() {
	s.stop <- nil
}

func (s *Service) revoke() error {

	_, err := s.client.Revoke(context.TODO(), s.leaseId)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("servide:%s stop\n", s.Name)
	return err
}
