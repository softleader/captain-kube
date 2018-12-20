package dockerctl

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"io"
	"log"
)

func Pull(image chart.Image) (io.ReadCloser, error) {
	ctx := context.Background()

	// Use DOCKER_HOST to set the url to the docker server.
	// Use DOCKER_API_VERSION to set the version of the API to reach, leave empty for latest.
	// Use DOCKER_CERT_PATH to load the TLS certificates from.
	// Use DOCKER_TLS_VERIFY to enable or disable TLS verification, off by default.
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Println("client init failed: ", err.Error())
		return nil, err
	}

	log.Println("pulling image: ", image)
	out, err := cli.ImagePull(ctx, image.String(), types.ImagePullOptions{})
	if err != nil {
		log.Println("pull image failed: ", image)
		return nil, err
	}

	return out, nil
}

func Retage(source chart.Image, target chart.Image) (io.ReadCloser, error) {
	ctx := context.Background()

	// Use DOCKER_HOST to set the url to the docker server.
	// Use DOCKER_API_VERSION to set the version of the API to reach, leave empty for latest.
	// Use DOCKER_CERT_PATH to load the TLS certificates from.
	// Use DOCKER_TLS_VERIFY to enable or disable TLS verification, off by default.
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Println("client init failed: ", err.Error())
		return nil, err
	}

	log.Println("taging image: ", source, " to ", target)
	if err := cli.ImageTag(ctx, source.String(), target.String()); err != nil {
		log.Println("tag image failed: ", source, " to ", target)
		return nil, err
	}

	log.Println("pushing image: ", target)
	out, err := cli.ImagePush(ctx, target.String(), types.ImagePushOptions{})
	if err := cli.ImageTag(ctx, source.String(), target.String()); err != nil {
		log.Println("push image failed: ", target)
		return nil, err
	}

	return out, nil
}
