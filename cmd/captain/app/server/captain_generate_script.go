package server

import (
	"context"
	"github.com/softleader/captain-kube/pkg/proto"
)

func (s *CaptainServer) GenerateScript(context.Context, *proto.GenerateScriptRequest) (*proto.GenerateScriptResponse, error) {
	return &proto.GenerateScriptResponse{
		Out: "Hello world!",
	}, nil
}
