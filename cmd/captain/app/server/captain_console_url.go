package server

import (
	"context"
	"fmt"
	"github.com/softleader/captain-kube/pkg/proto"
)

func (s *CaptainServer) ConsoleURL(ctx context.Context, req *captainkube_v2.ConsoleURLRequest) (*captainkube_v2.ConsoleURLResponse, error) {
	k8s, err := s.k8s()
	if err != nil {
		return nil, err
	}
	resp := &captainkube_v2.ConsoleURLResponse{
		Vendor: k8s.Server.GitCommit,
	}
	if k8s.Server.IsICP() {
		resp.Url = fmt.Sprintf("https://%s:%v", req.GetHost(), 8443)
		return resp, nil
	}
	if k8s.Server.IsGCP() {
		return nil, fmt.Errorf("GCP is not supported yet")
	}
	return nil, fmt.Errorf("unsupported kubernetes vendor: %v", k8s.Server.GitCommit)
}
