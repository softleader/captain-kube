package server

import (
	"github.com/softleader/captain-kube/pkg/ver"
	"os"
)

type capletServer struct {
	metadata *ver.BuildMetadata
	hostname string
}

func NewCapletServer(metadata *ver.BuildMetadata) (s *capletServer) {
	s = &capletServer{
		metadata: metadata,
	}
	s.hostname, _ = os.Hostname()
	return
}
