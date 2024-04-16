package kubelet

import (
	"context"
	"errors"
	"github.com/containerd/containerd"
	logger "github.com/sirupsen/logrus"
	core "minik8s/pkgs/apiobject"
)

type ImageController struct{}

func (c ImageController) CreateImage(client *containerd.Client, imageName string, pullPolicy core.ImagePullPolicy) (containerd.Image, error) {
	if client == nil {
		return nil, errors.New("no Client available")
	}
	ctx := context.Background()
	switch pullPolicy {
	case core.PullAlways:
		image, err := client.Pull(ctx, imageName, containerd.WithPullUnpack)
		if err != nil {
			logger.Errorf("pull image error: %s", err.Error())
		}
		return image, nil

	case core.PullIfNeeds:
		image, err := client.ImageService().Get(ctx, imageName)
		if err == nil {
			return containerd.NewImage(client, image), nil
		}
		containerdImage, err := client.Pull(ctx, imageName, containerd.WithPullUnpack)
		if err != nil {
			logger.Errorf("pull image error: %s", err.Error())
		}
		return containerdImage, nil

	case core.PullNever:
		image, err := client.ImageService().Get(ctx, imageName)
		if err != nil {
			return nil, err
		}
		return containerd.NewImage(client, image), nil
	}
	return nil, errors.New("unknown error")
}
