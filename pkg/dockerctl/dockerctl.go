package dockerctl

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/verbose"
	"io"
)

func Pull(out io.Writer, image chart.Image, registryAuth *proto.RegistryAuth) (io.ReadCloser, error) {
	ctx := context.Background()

	// Use DOCKER_HOST to set the url to the docker server.
	// Use DOCKER_API_VERSION to set the version of the API to reach, leave empty for latest.
	// Use DOCKER_CERT_PATH to load the TLS certificates from.
	// Use DOCKER_TLS_VERIFY to enable or disable TLS verification, off by default.
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	verbose.Fprintf(out, "pulling image: %s\n", image)
	opt := types.ImagePullOptions{}
	if registryAuth != nil {
		if opt.RegistryAuth, err = encode(registryAuth); err != nil {
			return nil, err
		}
	}
	rc, err := cli.ImagePull(ctx, image.String(), opt)
	if err != nil {
		return nil, err
	}

	return rc, nil
}

func ReTag(out io.Writer, source chart.Image, target chart.Image, registryAuth *proto.RegistryAuth) (io.ReadCloser, error) {
	ctx := context.Background()

	// Use DOCKER_HOST to set the url to the docker server.
	// Use DOCKER_API_VERSION to set the version of the API to reach, leave empty for latest.
	// Use DOCKER_CERT_PATH to load the TLS certificates from.
	// Use DOCKER_TLS_VERIFY to enable or disable TLS verification, off by default.
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	verbose.Fprintf(out, "tagging image from %q to %q \n", source, target)
	if err := cli.ImageTag(ctx, source.String(), target.String()); err != nil {
		return nil, err
	}

	verbose.Fprintf(out, "pushing image: %s \n", target)
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

func encode(ra *proto.RegistryAuth) (string, error) {
	b, err := json.Marshal(ra)
	if err != nil {
		return "", nil
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
