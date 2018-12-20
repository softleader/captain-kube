package dockerctl

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/proto"
	"io"
	"log"
)

func Pull(image chart.Image, registryAuth *proto.RegistryAuth) (io.ReadCloser, error) {
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
	opt := types.ImagePullOptions{}
	if registryAuth != nil {
		if opt.RegistryAuth, err = encode(registryAuth); err != nil {
			return nil, err
		}
	}
	out, err := cli.ImagePull(ctx, image.String(), opt)
	if err != nil {
		log.Println("pull image failed: ", image)
		return nil, err
	}

	return out, nil
}

func ReTag(source chart.Image, target chart.Image, registryAuth *proto.RegistryAuth) (io.ReadCloser, error) {
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
	opt := types.ImagePushOptions{}
	if registryAuth != nil {
		if opt.RegistryAuth, err = encode(registryAuth); err != nil {
			return nil, err
		}
	}
	out, err := cli.ImagePush(ctx, target.String(), opt)
	if err := cli.ImageTag(ctx, source.String(), target.String()); err != nil {
		log.Println("push image failed: ", target)
		return nil, err
	}

	return out, nil
}

// proto.RegistryAuth gen 出來的 Username 跟 Password 的 u 跟 p 都是小寫 json, docker 要大寫, 所以只能再轉換一次
type auth struct {
	Username string `json:"Username,omitempty"`
	Password string `json:"Password,omitempty"`
}

func encode(ra *proto.RegistryAuth) (string, error) {
	b, err := json.Marshal(auth{
		Username: ra.Username,
		Password: ra.Password,
	})
	if err != nil {
		return "", nil
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
