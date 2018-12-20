package dockerctl

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
	"log"
)

func Pull(host string, repo string, tag string) (io.ReadCloser, error) {
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
	image := fmt.Sprintf("%s/%s:%s", host, repo, tag) // 要可以 supports 沒有 host 的 image
	log.Println("pulling image: ", image)
	out, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		log.Println("pull image failed: ", image)
		return nil, err
	}

	return out, nil
}
