package resources

import (
	"context"
	"errors"
	core "minik8s/pkgs/apiobject"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/remotes/docker"
	logger "github.com/sirupsen/logrus"
)

type ImageController struct{}

func (c ImageController) CreateImage(client *containerd.Client, imageName string, pullPolicy core.ImagePullPolicy) (containerd.Image, error) {
	if client == nil {
		return nil, errors.New("no Client available")
	}
	ctx := context.Background()
	resolver := docker.NewResolver(docker.ResolverOptions{
		PlainHTTP: true,
	})
	switch pullPolicy {
	case core.PullAlways:
		image, err := client.Pull(ctx, imageName, containerd.WithPullUnpack, containerd.WithResolver(resolver))
		if err != nil {
			logger.Errorf("pull image error: %s", err.Error())
		}
		logger.Infof("pulled image: %v", image)
		return image, nil

	case core.PullIfNeeds:
		image, err := client.ImageService().Get(ctx, imageName)
		if err == nil {
			return containerd.NewImage(client, image), nil
		}
		containerdImage, err := client.Pull(ctx, imageName, containerd.WithPullUnpack, containerd.WithResolver(resolver))
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
