package model

const (
	TaskStatus_Ready   = 1000
	TaskStatus_Waiting = 1001
	TaskStatus_Doing   = 1002
	TaskStatus_Finish  = 1003
	TaskStatus_Error   = 1004
)

type TaskInfo struct {
	TaskId   string
	TaskName string
	Config   TaskConfig
}

type TaskConfig struct {
	TaskType int
	BeginAt  int64
	Status   int
}
