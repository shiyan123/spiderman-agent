package service

import (
	"fmt"
	"os/exec"
	"spiderman-agent/common/model"
	"sync"

	"github.com/robfig/cron"
)

type TaskService struct {
	infoMap map[string]*model.TaskInfo
	cronMap map[string]*cron.Cron
}

var (
	taskServiceOnce sync.Once
	taskService     *TaskService
)

func GetAccountService() *TaskService {
	taskServiceOnce.Do(func() {
		infoMap := make(map[string]*model.TaskInfo, 0)
		cronMap := make(map[string]*cron.Cron, 0)
		taskService = &TaskService{
			infoMap: infoMap,
			cronMap: cronMap,
		}
	})

	return taskService
}

func (t *TaskService) init(task *model.TaskInfo) {
	t.infoMap[task.TaskId] = task
	t.cronMap[task.TaskId] = cron.New()
	t.start(task.TaskId)
}

func (t *TaskService) start(id string) {
	spec := t.infoMap[id].Config.CronStr
	//todo 根据类型判断运行次数 重试次数 自动重启等
	t.cronMap[id].AddFunc(spec, func() {
		cmd := t.genCmd(id)
		_, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			fmt.Println(fmt.Sprintf("Failed to execute command: %s", cmd))
		}
	})
	t.cronMap[id].Start()
	select {}
}

func (t *TaskService) genCmd(id string) string {
	path := t.infoMap[id].Config.Path
	name := t.infoMap[id].Config.ProgramName
	cmd := fmt.Sprintf("cd ~ && cd %s	&& ./%s >> %s.log 2>&1 &",
		path, name, name)
	return cmd
}

func (t *TaskService) stop(task *model.TaskInfo) {
	t.cronMap[task.TaskId].Stop()
}
