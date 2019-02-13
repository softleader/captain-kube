package dockerd

import (
	"context"
	"github.com/fsouza/go-dockerclient"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/helm/chart"
)

func Rmi(log *logrus.Logger, force, dryRun bool, images ...*chart.Image) error {
	ctx := context.Background()

	// Use DOCKER_HOST to set the url to the docker server.
	// Use DOCKER_API_VERSION to set the version of the API to reach, leave empty for latest.
	// Use DOCKER_CERT_PATH to load the TLS certificates from.
	// Use DOCKER_TLS_VERIFY to enable or disable TLS verification, off by default.
	cli, err := docker.NewClientFromEnv()
	if err != nil {
		return err
	}

	// 參數準備
	options := docker.RemoveImageOptions{
		Context: ctx,
		Force:   force,
	}

	for _, i := range images {
		image := i.String()
		log.Infof("removing %s", image)
		if !dryRun {
			if err := cli.RemoveImageExtended(image, options); err != nil {
				return err
			}
		}
	}
	return nil
}
