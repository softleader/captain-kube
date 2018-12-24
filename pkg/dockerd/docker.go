package dockerd

import (
	"encoding/base64"
	"encoding/json"
	"github.com/softleader/captain-kube/pkg/proto"
)

func encode(ra *proto.RegistryAuth) (string, error) {
	b, err := json.Marshal(ra)
	if err != nil {
		return "", nil
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
