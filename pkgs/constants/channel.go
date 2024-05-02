package constants

import "fmt"

const (
	ChannelNode    string = "NODE"
	ChannelPod     string = "POD"
	ChannelService string = "SERVICE"
	ChannelReplica string = "REPLICASET"
)

const (
	ChannelCreate string = "CREATE"
	ChannelUpdate string = "UPDATE"
	ChannelDelete string = "DELETE"
)

var Channels = []string{ChannelNode, ChannelPod, ChannelService, ChannelReplica}
var Operations = []string{ChannelCreate, ChannelUpdate, ChannelDelete}

func GenerateChannelName(object string, chanType string) string {
	return fmt.Sprintf("%s-%s", object, chanType)
}
