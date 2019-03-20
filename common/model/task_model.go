package model

const (
	TaskStatus_Ready   = 1000
	TaskStatus_Waiting = 1001
	TaskStatus_Start   = 1002
	TaskStatus_Doing   = 1003
	TaskStatus_Finish  = 1004
	TaskStatus_Error   = 1005
	TaskStatus_Stop    = 1006
)

type TaskInfo struct {
	TaskId   string      `json:"taskId"`
	TaskName string      `json:"taskName"`
	Config   *TaskConfig `json:"config"`
}

type TaskConfig struct {
	TaskType        int    `json:"taskType"`        //任务类型
	CreatedAt       int64  `json:"createdAt"`       //创建时间
	LastBeginAt     int64  `json:"lastBeginAt"`     //上一次开始时间
	LastEndAt       int64  `json:"lastEndAt"`       //上一次结束时间
	ProgramUpdateAt int64  `json:"programUpdateAt"` //任务程序更新时间
	ProgramName     string `json:"programName"`     //运行程序
	CronStr         string `json:"cronStr"`         //cron 表达式
	Path            string `json:"path"`            //运行目录
	LogPath         string `json:"logPath"`         //日志目录
	Status          int    `json:"status"`          //运行状态
	Retry           int    `json:"retry"`           //重试次数
}
