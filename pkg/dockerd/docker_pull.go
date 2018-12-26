package dockerd

import (
	"context"
	"github.com/fsouza/go-dockerclient"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/proto"
)

func Pull(log *logrus.Logger, image chart.Image, registryAuth *proto.RegistryAuth) (err error) {
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

	// 參數準備
	options := docker.PullImageOptions{
		Context:      ctx,
		Tag:          image.Tag,
		Repository:   image.HostRepo(),
		OutputStream: log.Writer(),
	}

	// 第一此採用沒有帳密的方式, 若失敗則重試第二次, 第二次採用帳密
	if err := cli.PullImage(options, docker.AuthConfiguration{}); isDockerUnauthorized(err) && registryAuth != nil {
		err = cli.PullImage(options, docker.AuthConfiguration{
			Username: registryAuth.Username,
			Password: registryAuth.Password,
		})
	}

	return err
}
