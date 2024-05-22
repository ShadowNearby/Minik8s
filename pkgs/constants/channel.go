package constants

import "fmt"

const (
	ChannelNode    string = "NODE"
	ChannelPod     string = "POD"
	ChannelService string = "SERVICE"
	ChannelReplica string = "REPLICASET"
	ChannelHPA     string = "HPA"
)

const (
	ChannelCreate string = "CREATE"
	ChannelUpdate string = "UPDATE"
	ChannelDelete string = "DELETE"
)

var Channels = []string{ChannelNode, ChannelPod, ChannelService, ChannelReplica, ChannelHPA}
var Operations = []string{ChannelCreate, ChannelUpdate, ChannelDelete}
var OtherChannels = []string{ChannelPodSchedule}

func GenerateChannelName(object string, chanType string) string {
	return fmt.Sprintf("%s-%s", object, chanType)
}

// Additional Channels

const (
	ChannelPodSchedule string = "POD-SCHEDULE"
)
