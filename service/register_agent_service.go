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

type Service struct {
	Name      string
	Node      *model.Node
	stop      chan error
	leaseId   clientv3.LeaseID
	client    *clientv3.Client
	clientTTL int64
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
		Name:      genName(),
		Node:      model.InitNode(genID()),
		stop:      make(chan error),
		client:    cli,
		clientTTL: getClientTTL(),
	}

	// register service to etcd
	leaseId, err := s.register()
	if err != nil {
		return nil, err
	}

	//keepAlive service
	go func() {
		err = s.StartKeepAlive(leaseId)
		fmt.Printf("keepalive is err :%s  \n", err.Error())
	}()
	return s, nil
}

func (s *Service) register() (leaseId clientv3.LeaseID, err error) {

	value, err := s.getOriginalInfo()
	if value == "" || err != nil {
		return leaseId, err
	}

	resp, err := s.client.Grant(context.TODO(), s.clientTTL)
	if err != nil {
		return
	}
	_, err = s.client.Put(context.TODO(), s.Name, value, clientv3.WithLease(resp.ID))
	if err != nil {
		return
	}
	return resp.ID, nil
}

func (s *Service) getOriginalInfo() (string, error) {

	//get original service info
	fmt.Println("get original info")
	resp, err := s.client.Get(context.Background(), s.Name, clientv3.WithPrefix())
	if err != nil {
		return "", err
	}
	if len(resp.Kvs) > 0 {
		value := ""
		for _, v := range resp.Kvs {
			value = string(v.Value)
		}
		var n model.Node
		err = json.Unmarshal([]byte(value), &n)
		if err != nil {
			return "", err
		}
		s.Node = &n
	}
	body, err := json.Marshal(s.Node)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (s *Service) StartKeepAlive(leaseId clientv3.LeaseID) error {

	ch, err := s.client.KeepAlive(context.TODO(), leaseId)
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
		case _, ok := <-ch:
			if !ok {
				s.revoke()
				return errors.New("keep alive channel closed")
			} else {
				leaseId, err := s.keepAlive()
				if err != nil {
					return errors.New(fmt.Sprintf("leaseId: %d, keep alive send error", leaseId))
				}

				body, _ := json.Marshal(s.Node)
				fmt.Println("当前info：")
				fmt.Println(string(body))
			}
		}
	}
}

func (s *Service) keepAlive() (leaseId clientv3.LeaseID, err error) {

	MachineInfo, err := health.GetMachineHealth()
	if err != nil {
		return
	}
	s.Node.MachineInfo = MachineInfo

	body, _ := json.Marshal(s.Node)
	resp, err := s.client.Grant(context.TODO(), s.clientTTL)
	if err != nil {
		return
	}
	_, err = s.client.Put(context.TODO(), s.Name, string(body), clientv3.WithLease(resp.ID))
	if err != nil {
		return
	}
	return resp.ID, nil
}

func (s *Service) Stop() {
	s.stop <- errors.New("server stopped")
}

func (s *Service) revoke() error {

	_, err := s.client.Revoke(context.TODO(), s.leaseId)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("servide:%s stop\n", s.Name)
	return err
}

func genName() string {
	temp := fmt.Sprintf("%s-%s",
		app.GetApp().Config.Server.RegisterName,
		app.GetApp().Config.Server.RegisterIp)

	return fmt.Sprintf("services/%s/%s",
		app.GetApp().Config.Server.RegisterName,
		utils.EncodeStr(temp))
}

func genID() string {
	temp := fmt.Sprintf("%s-%s",
		app.GetApp().Config.Server.RegisterName,
		app.GetApp().Config.Server.RegisterIp)

	return utils.EncodeStr(temp)
}

func getClientTTL() int64 {
	if app.GetApp().Config.Server.ClientTTL == 0 {
		return 5
	}
	return app.GetApp().Config.Server.ClientTTL
}
