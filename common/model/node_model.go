package model

import (
	"spiderman-agent/app"
	"spiderman-agent/common/health"
)

type Node struct {
	ID          string                `json:"id"`
	IP          string                `json:"ip"`
	Name        string                `json:"name"`
	MachineInfo *health.MachineHealth `json:"machineInfo"`
	TaskMap     map[string]*TaskInfo  `json:"taskMap"`
}

func InitNode(id string) *Node {
	//obtain machine health
	machineInfo, err := health.GetMachineHealth()
	if err != nil {
		return nil
	}
	taskMap := make(map[string]*TaskInfo, 0)
	return &Node{
		ID:          id,
		IP:          app.GetApp().Config.Server.RegisterIp,
		Name:        app.GetApp().Config.Server.RegisterName,
		MachineInfo: machineInfo,
		TaskMap:     taskMap,
	}
}

func (n *Node) GetNodeInfo() *Node {
	return n
}
