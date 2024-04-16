package utils

import (
	"fmt"
	core "minik8s/pkgs/apiobject"
)

//type ContainerWalker struct {
//	Client  containerd.Client
//	OnFound Found
//}
//
//type Found struct {
//	Container  containerd.Container
//	Req        string // The raw request string. name, short ID, or long ID.
//	MatchIndex int    // Begins with 0, up to MatchCount - 1.
//	MatchCount int    // 1 on exact match. > 1 on ambiguous match. Never be <= 0.
//}
//
//type OnFound func(ctx context.Context, found Found) error

func GenerateContainerSpec(pConfig core.Pod, cConfig core.Container, args ...string) core.ContainerdSpec {
	var cSpec = core.ContainerdSpec{
		Namespace:    pConfig.MetaData.NameSpace,
		Image:        cConfig.Image,
		ID:           GenerateContainerIDByName(cConfig.Name),
		Name:         GenerateContainerName(pConfig, cConfig),
		VolumeMounts: generateVolMountsMap(cConfig.VolumeMounts),
		Cmd:          cConfig.Cmd,
		Envs:         generateEnvList(cConfig.Env),
		Resource:     cConfig.Resources.Limit,
		Labels:       pConfig.MetaData.Labels,
		PullPolicy:   cConfig.ImagePullPolicy,
	}
	if len(args) > 0 && len(args[0]) > 0 {
		linuxNamespace := args[0]
		var (
			ns = map[string]string{
				"cgroup":  linuxNamespace + "cgroup",
				"uts":     linuxNamespace + "uts",
				"network": linuxNamespace + "net",
				"ipc":     linuxNamespace + "ipc",
			}
		)
		cSpec.LinuxNamespace = ns
	}
	return cSpec
}

func GenerateContainerName(pConfig core.Pod, cConfig core.Container) string {
	return fmt.Sprintf("%s_%s", cConfig.Name, pConfig.MetaData.UUID)
}
