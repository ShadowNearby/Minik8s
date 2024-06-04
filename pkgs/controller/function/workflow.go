package function

import (
	core "minik8s/pkgs/apiobject"
	"minik8s/pkgs/apiserver/storage"
	constants "minik8s/pkgs/constants"
	worker "minik8s/pkgs/serverless/workflow"
	"minik8s/utils"

	logger "github.com/sirupsen/logrus"
)

type WorkFlowController struct{}

func (w *WorkFlowController) StartController() {
	triggerMessage := storage.RedisInstance.SubscribeChannel(constants.ChannelWorkflowTrigger)
	go func() {
		for message := range triggerMessage {
			err := w.HandleTrigger(message.Payload)
			if err != nil {
				logger.Errorf("handle trigger error: %s", err.Error())
			}
		}
	}()
	select {}
}

func (w *WorkFlowController) GetChannel() string {
	return constants.ChannelWorkflow
}

func (w *WorkFlowController) HandleCreate(message string) error {
	return nil
}

func (w *WorkFlowController) HandleUpdate(message string) error {
	return nil
}

func (w *WorkFlowController) HandleDelete(message string) error {
	return nil
}

func (w *WorkFlowController) HandleTrigger(message string) error {
	var request core.WorkFlowTriggerRequest
	err := utils.JsonUnMarshal(message, &request)
	if err != nil {
		logger.Error("error parse trigger request")
		return err
	}
	logger.Infof("trigger workflow %s", request.Name)
	var workflow core.Workflow
	wfTxt := utils.GetObjectWONamespace(core.ObjWorkflow, request.Name)
	err = utils.JsonUnMarshal(wfTxt, &workflow)
	workflow.States = utils.ParseWorkStateMap(workflow.States) // parse again
	if err != nil {
		logger.Error("error parse workflow")
		return err
	}
	result, err := worker.ExecuteWorkflow(&workflow, []byte(request.Params))
	if err != nil {
		logger.Errorf("execute workflow error: %s", err.Error())
		return nil
	}
	logger.Infof("workflow execute result: %s", string(result))
	return nil
}
