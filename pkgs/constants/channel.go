package constants

import "fmt"

const (
	ChannelNode     string = "NODE"
	ChannelPod      string = "POD"
	ChannelService  string = "SERVICE"
	ChannelReplica  string = "REPLICASET"
	ChannelFunction string = "FUNCTION"
	ChannelHPA      string = "HPA"
)

const (
	ChannelCreate  string = "CREATE"
	ChannelUpdate  string = "UPDATE"
	ChannelDelete  string = "DELETE"
	ChannelTrigger string = "TRIGGER"
)

var Channels = []string{ChannelNode, ChannelPod, ChannelService, ChannelReplica, ChannelHPA, ChannelFunction}
var Operations = []string{ChannelCreate, ChannelUpdate, ChannelDelete, ChannelTrigger}
var OtherChannels = []string{ChannelPodSchedule}

func GenerateChannelName(object string, chanType string) string {
	return fmt.Sprintf("%s-%s", object, chanType)
}

// Additional Channels

const (
	ChannelPodSchedule string = "POD-SCHEDULE"
)
