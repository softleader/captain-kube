package dockerctl

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
	"log"
	"os"
)

func Pull(host string, repo string, tag string) (io.ReadCloser, error) {
	// docker client
	ctx := context.Background()

	cliVsn, exist := os.LookupEnv("DOCKER_CLIENT_VERSION")
	if !exist {
		cliVsn = "1.39"
	}

	cli, err := client.NewClientWithOpts(client.WithVersion(cliVsn))
	if err != nil {
		log.Println("client init failed, version: ", cliVsn)
		return nil, err
	}

	// docker pull image
	image := host + "/" + repo
	log.Println("pulling image: ", image)
	out, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		log.Println("pull image failed: ", image)
		return nil, err
	}
	defer out.Close()

	return out, nil
}
