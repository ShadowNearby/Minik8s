package constants

import "fmt"

const (
	ChannelNode     string = "NODE"
	ChannelPod      string = "POD"
	ChannelService  string = "SERVICE"
	ChannelReplica  string = "REPLICASET"
	ChannelFunction string = "FUNCTION"
	ChannelHPA      string = "HPA"
	ChannelEndpoint string = "ENDPOINT"
	ChannelTask     string = "TASK"
	ChannelWorkflow string = "WORKFLOW"
)

const (
	ChannelCreate string = "CREATE"
	ChannelUpdate string = "UPDATE"
	ChannelDelete string = "DELETE"
)

var Channels = []string{ChannelNode, ChannelPod, ChannelService, ChannelReplica, ChannelHPA, ChannelFunction, ChannelTask, ChannelWorkflow}
var Operations = []string{ChannelCreate, ChannelUpdate, ChannelDelete}
var OtherChannels = []string{ChannelPodSchedule, ChannelFunctionTrigger, ChannelWorkflowTrigger}

func GenerateChannelName(object string, chanType string) string {
	return fmt.Sprintf("%s-%s", object, chanType)
}

// Additional Channels

const (
	ChannelPodSchedule     string = "POD-SCHEDULE"
	ChannelFunctionTrigger string = "FUNCTION-TRIGGER"
	ChannelWorkflowTrigger string = "WORKFLOW-TRIGGER"
)
