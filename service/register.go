package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"log"
	"spiderman-agent/app"
	"spiderman-agent/utils"
	"time"
)

type ServicesInfo struct {
	IP   string
	Name string
}

type Service struct {
	Name    string
	Info    *ServicesInfo
	stop    chan error
	leaseId clientv3.LeaseID
	client  *clientv3.Client
}

func RegisterService() (*Service, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   app.GetApp().Config.CenterServer.Ips,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	return &Service{
		Name:   fmt.Sprintf("services/%s/%s", app.GetApp().Config.Server.RegisterName, utils.EncodeTimeStr(time.Now())),
		Info:   initServicesInfo(),
		stop:   make(chan error),
		client: cli,
	}, err
}

func initServicesInfo() *ServicesInfo {
	return &ServicesInfo{
		IP:   app.GetApp().Config.Server.RegisterIp,
		Name: app.GetApp().Config.Server.RegisterName,
	}
}

func (s *Service) Start() error {

	ch, err := s.keepAlive()
	if err != nil {
		return err
	}

	for {
		select {
		case err := <-s.stop:
			s.revoke()
			return err
		case <-s.client.Ctx().Done():
			return errors.New("server closed")
		case ka, ok := <-ch:
			if !ok {
				log.Println("keep alive channel closed")
				s.revoke()
				return nil
			} else {
				log.Printf("Recv reply from service: %s, ttl:%d", s.Name, ka.TTL)
			}
		}
	}
}

func (s *Service) Stop() {
	s.stop <- nil
}

func (s *Service) keepAlive() (<-chan *clientv3.LeaseKeepAliveResponse, error) {

	info := &s.Info
	value, err := json.Marshal(info)

	if err != nil {
		return nil, err
	}
	resp, err := s.client.Grant(context.TODO(), 5)
	if err != nil {
		return nil, err
	}
	_, err = s.client.Put(context.TODO(), s.Name, string(value), clientv3.WithLease(resp.ID))
	if err != nil {
		return nil, err
	}
	s.leaseId = resp.ID

	return s.client.KeepAlive(context.TODO(), resp.ID)
}

func (s *Service) revoke() error {

	_, err := s.client.Revoke(context.TODO(), s.leaseId)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("servide:%s stop\n", s.Name)
	return err
}
