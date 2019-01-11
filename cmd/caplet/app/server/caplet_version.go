package server

import (
	"context"
	"github.com/softleader/captain-kube/pkg/proto"
)

func (s *CapletServer) Version(ctx context.Context, req *captainkube_v2.VersionRequest) (resp *captainkube_v2.VersionResponse, err error) {
	resp = &captainkube_v2.VersionResponse{
		Hostname: s.hostname,
	}
	if req.GetFull() {
		resp.Version = s.metadata.FullString()
	} else {
		resp.Version = s.metadata.String()
	}
	return
}
