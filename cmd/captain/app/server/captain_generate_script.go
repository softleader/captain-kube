package server

import (
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/sio"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func (s *CaptainServer) GenerateScript(req *proto.GenerateScriptRequest, stream proto.Captain_GenerateScriptServer) error {
	tmp, err := ioutil.TempDir(os.TempDir(), "generate-script-")
	if err != nil {
		return err
	}
	chartPath := filepath.Join(tmp, req.GetChart().GetFileName())
	if err := ioutil.WriteFile(chartPath, req.GetChart().GetContent(), 0644); err != nil {
		return err
	}

	sout := sio.NewStreamWriter(func(p []byte) error {
		return stream.Send(&proto.GenerateScriptResponse{
			Msg: p,
		})
	})

	tpls, err := chart.LoadArchive(s.out, chartPath)
	if err != nil {
		return err
	}

	if from, to := strings.TrimSpace(req.GetRetag().GetFrom()), strings.TrimSpace(req.GetRetag().GetTo()); from != "" && to != "" {
		if err := tpls.GenerateReTagScript(sout, from, to); err != nil {
			return err
		}
	}

	if req.Pull {
		if err := tpls.GeneratePullScript(sout); err != nil {
			return err
		}
	}

	if req.Load {
		if err := tpls.GenerateLoadScript(sout); err != nil {
			return err
		}
	}

	if req.Save {
		if err := tpls.GenerateSaveScript(sout); err != nil {
			return err
		}
	}

	return nil
}
