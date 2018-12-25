package dockerd

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/proto"
	"io"
)

func ReTag(log *logrus.Logger, source chart.Image, target chart.Image, registryAuth *proto.RegistryAuth) (io.ReadCloser, error) {
	ctx := context.Background()

	// Use DOCKER_HOST to set the url to the docker server.
	// Use DOCKER_API_VERSION to set the version of the API to reach, leave empty for latest.
	// Use DOCKER_CERT_PATH to load the TLS certificates from.
	// Use DOCKER_TLS_VERIFY to enable or disable TLS verification, off by default.
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	log.Printf("tagging image from %q to %q \n", source, target)
	if err := cli.ImageTag(ctx, source.String(), target.String()); err != nil {
		return nil, err
	}

	log.Printf("pushing image: %s \n", target)
	opt := types.ImagePushOptions{}
	if registryAuth != nil {
		if opt.RegistryAuth, err = encode(registryAuth); err != nil {
			return nil, err
		}
	}
	rc, err := cli.ImagePush(ctx, target.String(), opt)
	if err := cli.ImageTag(ctx, source.String(), target.String()); err != nil {
		return nil, err
	}

	return rc, nil
}
