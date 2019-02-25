package dockerd

import (
	"encoding/base64"
	"encoding/json"
	"github.com/fsouza/go-dockerclient"
	pb "github.com/softleader/captain-kube/pkg/proto"
)

func encode(ra *pb.RegistryAuth) (string, error) {
	b, err := json.Marshal(ra)
	if err != nil {
		return "", nil
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func isDockerUnauthorized(err error) bool {
	if err != nil {
		if _, ok := err.(*docker.Error); ok {
			// 從 header 的文字判斷是不好的做法, 畢竟 docker daemon 升級後也許會調整
			// 因此這邊先當成只要是 docker 回復的錯誤就是 unauthorized
			// return strings.Contains(derr.Message, "missing X-Registry-Auth")
			return true
		}
	}
	return false
}
