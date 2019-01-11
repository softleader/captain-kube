package dockerd

import (
	"context"
	"github.com/fsouza/go-dockerclient"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/proto"
)

func Pull(log *logrus.Logger, image chart.Image, registryAuth *captainkube_v2.RegistryAuth) error {
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
	options := docker.PullImageOptions{
		Context:      ctx,
		Tag:          image.Tag,
		Repository:   image.HostRepo(),
		OutputStream: log.Out,
	}

	// 第一此採用沒有帳密的方式, 若失敗且為驗證相關的問題, 則重試第二次並採用帳密
	if err := cli.PullImage(options, docker.AuthConfiguration{}); err != nil {
		if registryAuth == nil || !isDockerUnauthorized(err) {
			return err
		}
		// 為了避免 ineffectual assignments, so lets make a new var of error
		if erragain := cli.PullImage(options, docker.AuthConfiguration{
			Username: registryAuth.Username,
			Password: registryAuth.Password,
		}); erragain != nil {
			return erragain
		}
	}

	return nil
}

func PullFromTemplates(log *logrus.Logger, tpls chart.Templates, auth *captainkube_v2.RegistryAuth) error {
	for _, tpl := range tpls {
		for _, image := range tpl {
			log.Debugln("pulling image:", image.String())
			if err := Pull(log, *image, auth); err != nil {
				return err
			}
		}
	}
	return nil
}
