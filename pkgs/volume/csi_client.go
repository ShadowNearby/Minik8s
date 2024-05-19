package volume

import (
	"context"
	"minik8s/config"
	"net"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CsiClient struct {
	ControllerClient      csi.ControllerClient
	NodeClient            csi.NodeClient
	IdentityClient        csi.IdentityClient
	GroupControllerClient csi.GroupControllerClient
	Context               context.Context
}

var CsiClientInstance = NewCsiClient(config.CsiSockAddr)

func NewCsiClient(addr string) *CsiClient {
	network := "unix"
	connection, err := grpc.Dial(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, target string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, network, target)
		}))
	if err != nil {
		logrus.Errorf("error connecting to csi %s: %s", addr, err.Error())
		return nil
	}
	return &CsiClient{
		ControllerClient:      csi.NewControllerClient(connection),
		NodeClient:            csi.NewNodeClient(connection),
		IdentityClient:        csi.NewIdentityClient(connection),
		GroupControllerClient: csi.NewGroupControllerClient(connection),
		Context:               context.Background(),
	}
}
