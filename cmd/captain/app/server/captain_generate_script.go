package server

import (
	"bytes"
	"context"
	"github.com/softleader/captain-kube/pkg/arc"
	"github.com/softleader/captain-kube/pkg/helm"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/proto"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const template = "t"

func (s *CaptainServer) GenerateScript(c context.Context, req *proto.GenerateScriptRequest) (resp *proto.GenerateScriptResponse, err error) {
	path, err := ioutil.TempDir(os.TempDir(), "captain-generate-script")
	if err != nil {
		return
	}
	chartFile := "" // TODO
	chartPath := filepath.Join(path, req.GetChart().GetFileName())
	if err = arc.Extract(s.out, chartFile, chartPath); err != nil {
		return
	}
	tplPath := filepath.Join(chartPath, template)
	if err := helm.Template(s.out, chartPath, tplPath); err != nil {
		return
	}
	images, err := chart.CollectImages(tplPath)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	if from, to := strings.TrimSpace(req.GetRetag().GetFrom()), strings.TrimSpace(req.GetRetag().GetTo()); from != "" && to != "" {
		script, err := images.GenerateReTagScript(from, to)
		if err != nil {
			return nil, err
		}
		buf.Write(script.Bytes())
	}

	if req.Pull {
		script, err := images.GeneratePullScript()
		if err != nil {
			return nil, err
		}
		buf.Write(script.Bytes())
	}

	if req.Load {
		script, err :=  images.GenerateLoadScript()
		if err != nil {
			return nil, err
		}
		buf.Write(script.Bytes())
	}

	if req.Save {
		script, err := images.GenerateSaveScript()
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
