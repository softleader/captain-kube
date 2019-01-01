package server

import (
	"context"
	"github.com/softleader/captain-kube/pkg/proto"
)

func (s *CapletServer) Version(ctx context.Context, req *proto.VersionRequest) (resp *proto.VersionResponse, err error) {
	resp = &proto.VersionResponse{
		Hostname: s.hostname,
	}
	if req.GetLong() {
		resp.Version = s.metadata.LongString()
	} else {
		resp.Version = s.metadata.String()
	}
	return
}
