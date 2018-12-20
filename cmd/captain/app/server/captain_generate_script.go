package server

import (
	"bytes"
	"context"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/proto"
	"strings"
)

func (s *CaptainServer) GenerateScript(c context.Context, req *proto.GenerateScriptRequest) (resp *proto.GenerateScriptResponse, err error) {
	chartFile := req.GetChart().GetFileName() // TODO
	tpls, err := chart.LoadArchive(s.out, chartFile)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	if from, to := strings.TrimSpace(req.GetRetag().GetFrom()), strings.TrimSpace(req.GetRetag().GetTo()); from != "" && to != "" {
		script, err := tpls.GenerateReTagScript(from, to)
		if err != nil {
			return nil, err
		}
		buf.Write(script.Bytes())
	}

	if req.Pull {
		script, err := tpls.GeneratePullScript()
		if err != nil {
			return nil, err
		}
		buf.Write(script.Bytes())
	}

	if req.Load {
		script, err := tpls.GenerateLoadScript()
		if err != nil {
			return nil, err
		}
		buf.Write(script.Bytes())
	}

	if req.Save {
		script, err := tpls.GenerateSaveScript()
		if err != nil {
			return nil, err
		}
		buf.Write(script.Bytes())
	}

	resp = &proto.GenerateScriptResponse{
		Out: buf.String(),
	}

	return
}
