package server

import (
	"github.com/softleader/captain-kube/pkg/release"
	"os"
)

type CapletServer struct {
	metadata *release.Metadata
	hostname string
}

func NewCapletServer(metadata *release.Metadata) (s *CapletServer) {
	s = &CapletServer{
		metadata: metadata,
	}
	s.hostname, _ = os.Hostname()
	return
}
