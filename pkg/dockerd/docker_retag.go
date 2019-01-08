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

	// 參數準備
	options := docker.PushImageOptions{
		Context:      ctx,
		Name:         target.HostRepo(),
		Tag:          target.Tag,
		OutputStream: log.Out,
	}

	// 第一此採用沒有帳密的方式, 若失敗且為驗證相關的問題, 則重試第二次並採用帳密
	if err := cli.PushImage(options, docker.AuthConfiguration{}); err != nil {
		if registryAuth == nil || !isDockerUnauthorized(err) {
			return err
		}
		// 為了避免 ineffectual assignments, so lets make a new var of error
		if erragain := cli.PushImage(options, docker.AuthConfiguration{
			Username: registryAuth.Username,
			Password: registryAuth.Password,
		}); erragain != nil {
			return erragain
		}
	}

	return nil
}

func ReTagFromTemplates(log *logrus.Logger, tpls chart.Templates, retag *proto.ReTag, auth *proto.RegistryAuth) error {
	for _, tpl := range tpls {
		for _, image := range tpl {
			if image.Host == retag.From {
				log.Println("syncing ", image)
				if err := ReTag(log, *image, chart.Image{
					Host: retag.To,
					Repo: image.Repo,
					Tag:  image.Tag,
				}, auth); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
