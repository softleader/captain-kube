package server

import (
	"context"
	"github.com/softleader/captain-kube/pkg/proto"
)

func (s *capletServer) Version(ctx context.Context, req *proto.VersionRequest) (*proto.VersionResponse, error) {
	return &proto.VersionResponse{
		Hostname: s.hostname,
		Version:  s.metadata.String(req.GetShort()),
	}, nil
}
