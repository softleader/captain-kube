package server

import (
	"github.com/softleader/captain-kube/pkg/ver"
	"os"
)

type CapletServer struct {
	metadata *ver.BuildMetadata
	hostname string
}

func NewCapletServer(metadata *ver.BuildMetadata) (s *CapletServer) {
	s = &CapletServer{
		metadata: metadata,
	}
	s.hostname, _ = os.Hostname()
	return
}
