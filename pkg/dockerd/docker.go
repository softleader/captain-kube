package dockerd

import (
	"encoding/base64"
	"encoding/json"
	"github.com/fsouza/go-dockerclient"
	"github.com/softleader/captain-kube/pkg/proto"
)

func encode(ra *proto.RegistryAuth) (string, error) {
	b, err := json.Marshal(ra)
	if err != nil {
		return "", nil
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func isDockerUnauthorized(err error) bool {
	if err != nil {
		if _, ok := err.(*docker.Error); ok {
			return true
			// return derr.Status == 500 && strings.HasSuffix(derr.Message, "no basic auth credentials")
		}
	}
	return false
}
