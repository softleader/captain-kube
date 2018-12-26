package dockerd

import (
	"context"
	"github.com/fsouza/go-dockerclient"
	"github.com/sirupsen/logrus"
)

func Prune(log *logrus.Logger) error {
	ctx := context.Background()
	cli, err := docker.NewClientFromEnv()
	if err != nil {
		return err
	}

	if cp, err := cli.PruneContainers(docker.PruneContainersOptions{Context: ctx}); err != nil {
		return err
	} else {
		log.Info("deleted containers:")
		for _, del := range cp.ContainersDeleted {
			log.Info(del)
		}
		log.Infof("total reclaimed space: %v", cp.SpaceReclaimed)
	}

	if ip, err := cli.PruneImages(docker.PruneImagesOptions{Context: ctx}); err != nil {
		return err
	} else {
		log.Info("deleted images:")
		for _, del := range ip.ImagesDeleted {
			log.Info("deleted:", del.Deleted)
			log.Info("untagged:", del.Untagged)
		}
		log.Infof("total reclaimed space: %v", ip.SpaceReclaimed)
	}

	if np, err := cli.PruneNetworks(docker.PruneNetworksOptions{Context: ctx}); err != nil {
		return err
	} else {
		log.Info("deleted networks:")
		for _, del := range np.NetworksDeleted {
			log.Info(del)
		}
	}

	if vp, err := cli.PruneVolumes(docker.PruneVolumesOptions{Context: ctx}); err != nil {
		return err
	} else {
		log.Info("deleted volumes:")
		for _, del := range vp.VolumesDeleted {
			log.Info(del)
		}
	}
	return nil
}
