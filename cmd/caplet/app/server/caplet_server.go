package server

import (
	"github.com/softleader/captain-kube/pkg/version"
	"os"
)

type CapletServer struct {
	metadata *version.BuildMetadata
	hostname string
}

func NewCapletServer(metadata *version.BuildMetadata) (s *CapletServer) {
	s = &CapletServer{
		metadata: metadata,
	}
	s.hostname, _ = os.Hostname()
	return
}
