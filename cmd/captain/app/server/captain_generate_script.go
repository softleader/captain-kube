package server

import (
	"bytes"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/proto"
	"strings"
)

func (s *CaptainServer) GenerateScript(req *proto.GenerateScriptRequest, stream proto.Captain_GenerateScriptServer) error {
	chartFile := req.GetChart().GetFileName() // TODO
	tpls, err := chart.LoadArchive(s.out, chartFile)
	if err != nil {
		return err
	}

	var buf bytes.Buffer

	if from, to := strings.TrimSpace(req.GetRetag().GetFrom()), strings.TrimSpace(req.GetRetag().GetTo()); from != "" && to != "" {
		script, err := tpls.GenerateReTagScript(from, to)
		if err != nil {
			return err
		}
		if err := stream.Send(&proto.GenerateScriptResponse{
			Out: script.Bytes(),
		}); err != nil {
			return err
		}
	}

	if req.Pull {
		script, err := tpls.GeneratePullScript()
		if err != nil {
			return err
		}
		if err := stream.Send(&proto.GenerateScriptResponse{
			Out: script.Bytes(),
		}); err != nil {
			return err
		}
	}

	if req.Load {
		script, err := tpls.GenerateLoadScript()
		if err != nil {
			return err
		}
		if err := stream.Send(&proto.GenerateScriptResponse{
			Out: script.Bytes(),
		}); err != nil {
			return err
		}
	}

	if req.Save {
		script, err := tpls.GenerateSaveScript()
		if err != nil {
			return err
		}
		if err := stream.Send(&proto.GenerateScriptResponse{
			Out: script.Bytes(),
		}); err != nil {
			return err
		}
	}

	return nil
}
