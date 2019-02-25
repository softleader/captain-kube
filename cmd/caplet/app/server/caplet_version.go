package server

import (
	"context"
	pb "github.com/softleader/captain-kube/pkg/proto"
)

// Version 回傳 caplet 的版本
func (s *CapletServer) Version(ctx context.Context, req *pb.VersionRequest) (resp *pb.VersionResponse, err error) {
	resp = &pb.VersionResponse{
		Hostname: s.hostname,
	}
	if req.GetFull() {
		resp.Version = s.metadata.FullString()
	} else {
		resp.Version = s.metadata.String()
	}
	return
}
