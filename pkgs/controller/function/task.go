package function

import (
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/constants"
	"minik8s/pkgs/serverless/activator"
	"minik8s/utils"
	"sync"

	"github.com/robfig/cron/v3"
	logger "github.com/sirupsen/logrus"
)

type TaskController struct {
	cronManager *cron.Cron
	entryIDMap  map[string]cron.EntryID
	Mutex       sync.Mutex
}

func (t *TaskController) StartTaskController() {
	t.cronManager = cron.New()
	t.entryIDMap = make(map[string]cron.EntryID)
	allTasks := getAllTasks()
	logger.Infof("old task length: %d", len(allTasks))
	for _, task := range allTasks {
		funcName := task.Spec.Sink.Ref.Name
		params := task.Spec.JsonData
		entryID, err := t.cronManager.AddFunc(task.Spec.Schedule, func() {
			logger.Infof("do task trigger func")
			_, err := activator.TriggerFunc(funcName, []byte(params))
			if err != nil {
				logger.Errorf("trigger function %s error", funcName)
				return
			}
		})
		if err != nil {
			logger.Errorf("register task error: %s", err.Error())
		}
		t.entryIDMap[funcName] = entryID
	}
	t.cronManager.Start()
}

func getAllTasks() []core.PingSource {
	tasksTxt := utils.GetObjectWONamespace(core.ObjFunction, string(core.ObjTask))
	var tasks []core.PingSource
	err := utils.JsonUnMarshal(tasksTxt, &tasks)
	if err != nil {
		logger.Errorf("there's no task now")
	}
	return tasks
}

func (t *TaskController) GetChannel() string {
	return constants.ChannelTask
}

func (t *TaskController) HandleCreate(message string) error {
	logger.Infof("handle create")
	var pingSource core.PingSource
	utils.JsonUnMarshal(message, &pingSource)
	funcName := pingSource.Spec.Sink.Ref.Name
	params := pingSource.Spec.JsonData

	entryID, err := t.cronManager.AddFunc(pingSource.Spec.Schedule, func() {
		result, err := activator.TriggerFunc(funcName, []byte(params))
		if err != nil {
			logger.Errorf("trigger function %s error", funcName)
			return
		}
		triggerResult := core.TriggerResult{
			ID:     pingSource.ID,
			Result: result,
		}
		utils.SaveTriggerResult(core.ObjFunction, triggerResult)
	})
	if err != nil {
		logger.Errorf("register task error: %s", err.Error())
		return err
	}
	t.Mutex.Lock()
	t.entryIDMap[funcName] = entryID
	t.Mutex.Unlock()
	return nil
}

func (t *TaskController) HandleUpdate(message string) error {
	var pingSource core.PingSource
	utils.JsonUnMarshal(message, &pingSource)
	funcName := pingSource.Spec.Sink.Ref.Name
	params := pingSource.Spec.JsonData
	// check if the task is under management
	if entryID, ok := t.entryIDMap[funcName]; ok {
		t.cronManager.Remove(entryID)
		entryID, err := t.cronManager.AddFunc(pingSource.Spec.Schedule, func() {
			_, err := activator.TriggerFunc(funcName, []byte(params))
			if err != nil {
				logger.Errorf("trigger function error: %s", err.Error())
				return
			}
		})
		if err != nil {
			logger.Errorf("register task error: %s", err.Error())
			return err
		}
		t.entryIDMap[funcName] = entryID
	}
	return nil
}

func (t *TaskController) HandleDelete(message string) error {
	var pingSource core.PingSource
	utils.JsonUnMarshal(message, &pingSource)
	funcName := pingSource.Spec.Sink.Ref.Name
	// check if the task is under management
	if entryID, ok := t.entryIDMap[funcName]; ok {
		t.cronManager.Remove(entryID)
	}
	return nil
}
