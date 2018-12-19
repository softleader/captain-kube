package dockerctl

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
	"log"
)

func Pull(host string, repo string, tag string) (io.ReadCloser, error) {
	// docker client
	ctx := context.Background()

	//cliVsn, exist := os.LookupEnv("DOCKER_CLIENT_VERSION")
	//if !exist {
	//	cliVsn = "1.39"
	//}

	// cli, err := client.NewClientWithOpts(client.WithVersion(cliVsn))
	// Use DOCKER_HOST to set the url to the docker server.
	// Use DOCKER_API_VERSION to set the version of the API to reach, leave empty for latest.
	// Use DOCKER_CERT_PATH to load the TLS certificates from.
	// Use DOCKER_TLS_VERIFY to enable or disable TLS verification, off by default.
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Println("client init failed: ", err.Error())
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

	return out, nil
}
