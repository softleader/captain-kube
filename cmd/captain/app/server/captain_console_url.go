package server

import (
	"context"
	"fmt"
	"github.com/softleader/captain-kube/pkg/proto"
)

func (s *CaptainServer) ConsoleURL(ctx context.Context, req *captainkube_v2.ConsoleURLRequest) (*captainkube_v2.ConsoleURLResponse, error) {
	resp := &captainkube_v2.ConsoleURLResponse{
		Vendor: s.K8s.ServerVersion.GitCommit,
	}
	if s.K8s.ServerVersion.IsICP() {
		resp.Url = fmt.Sprintf("https://%s:%v", req.GetHost(), 8443)
		return resp, nil
	}
	if s.K8s.ServerVersion.IsGCP() {
		return nil, fmt.Errorf("GCP is not supported yet")
	}
	return nil, fmt.Errorf("unsupported kubernetes vendor: %v", s.K8s.ServerVersion.GitCommit)
}
