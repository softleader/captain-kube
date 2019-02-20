package server

import (
	"context"
	"github.com/softleader/captain-kube/pkg/proto"
)

// Version 回傳 caplet 的版本
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
