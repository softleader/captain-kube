package dockerd

import (
	"context"
	"github.com/fsouza/go-dockerclient"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/proto"
)

func ReTag(log *logrus.Logger, source chart.Image, target chart.Image, registryAuth *proto.RegistryAuth) error {
	ctx := context.Background()

	// Use DOCKER_HOST to set the url to the docker server.
	// Use DOCKER_API_VERSION to set the version of the API to reach, leave empty for latest.
	// Use DOCKER_CERT_PATH to load the TLS certificates from.
	// Use DOCKER_TLS_VERIFY to enable or disable TLS verification, off by default.
	cli, err := docker.NewClientFromEnv()
	if err != nil {
		return err
	}

	log.Printf("tagging image from %q to %q \n", source.String(), target.String())
	if err := cli.TagImage(source.String(), docker.TagImageOptions{
		Context: ctx,
		Force:   true,
		Repo:    target.HostRepo(),
		Tag:     target.Tag,
	}); err != nil {
		return err
	}

	log.Printf("pushing image: %s \n", target.String())
	var auth docker.AuthConfiguration
	if registryAuth != nil {
		auth = docker.AuthConfiguration{
			Username: registryAuth.Username,
			Password: registryAuth.Password,
		}
	}
	if err := cli.PushImage(docker.PushImageOptions{
		Context:      ctx,
		Registry:     target.Host,
		Name:         target.Repo,
		Tag:          target.Tag,
		OutputStream: log.Writer(),
	}, auth); err != nil {
		return err
	}

	return nil
}
