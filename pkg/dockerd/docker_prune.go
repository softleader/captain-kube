package dockerd

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func Prune(log *logrus.Logger) error {
	ctx := context.Background()

	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	args := filters.NewArgs()

	if cp, err := cli.ContainersPrune(ctx, args); err != nil {
		return err
	} else {
		log.Info("deleted containers:")
		for _, del := range cp.ContainersDeleted {
			log.Info(del)
		}
		log.Infof("total reclaimed space: %v", cp.SpaceReclaimed)
	}

	if ip, err := cli.ImagesPrune(ctx, args); err != nil {
		return err
	} else {
		log.Info("deleted images:")
		for _, del := range ip.ImagesDeleted {
			log.Info(del)
		}
		log.Infof("total reclaimed space: %v", ip.SpaceReclaimed)
	}

	if np, err := cli.NetworksPrune(ctx, args); err != nil {
		return err
	} else {
		log.Info("deleted networks:")
		for _, del := range np.NetworksDeleted {
			log.Info(del)
		}
	}

	if vp, err := cli.VolumesPrune(ctx, args); err != nil {
		return err
	} else {
		log.Info("deleted volumes:")
		for _, del := range vp.VolumesDeleted {
			log.Info(del)
		}
	}
	return nil
}
