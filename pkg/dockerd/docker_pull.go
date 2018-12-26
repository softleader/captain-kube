package dockerd

import (
	"context"
	"github.com/fsouza/go-dockerclient"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/proto"
)

func Pull(log *logrus.Logger, image chart.Image, registryAuth *proto.RegistryAuth) error {
	ctx := context.Background()

	// Use DOCKER_HOST to set the url to the docker server.
	// Use DOCKER_API_VERSION to set the version of the API to reach, leave empty for latest.
	// Use DOCKER_CERT_PATH to load the TLS certificates from.
	// Use DOCKER_TLS_VERIFY to enable or disable TLS verification, off by default.
	cli, err := docker.NewClientFromEnv()
	if err != nil {
		return err
	}

	log.Printf("pulling image: %s\n", image.String())

	var auth docker.AuthConfiguration
	if registryAuth != nil {
		auth = docker.AuthConfiguration{
			Username: registryAuth.Username,
			Password: registryAuth.Password,
		}
	}

	if err := cli.PullImage(docker.PullImageOptions{
		Context:      ctx,
		Tag:          image.Tag,
		Repository:   image.HostRepo(),
		OutputStream: log.Writer(),
	}, auth); err != nil {
		return err
	}

	return nil
}
