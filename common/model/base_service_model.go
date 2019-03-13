package model

import (
	"spiderman-agent/app"
	"spiderman-agent/common/health"
)

type ServiceInfo struct {
	IP          string
	Name        string
	MachineInfo *health.MachineHealth
	TaskMap     map[string]*TaskInfo
}

func InitServiceInfo() *ServiceInfo {
	//obtain machine health
	machineInfo, err := health.GetMachineHealth()
	if err != nil {
		return nil
	}
	taskMap := make(map[string]*TaskInfo, 0)
	return &ServiceInfo{
		IP:          app.GetApp().Config.Server.RegisterIp,
		Name:        app.GetApp().Config.Server.RegisterName,
		MachineInfo: machineInfo,
		TaskMap:     taskMap,
	}
}

func (s *ServiceInfo) GetServiceInfo() *ServiceInfo {
	return s
}

