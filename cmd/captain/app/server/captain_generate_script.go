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
	images, err := chart.CollectImages(tplPath, func(image chart.Image) bool {
		return true
	}, func(image chart.Image) chart.Image {
		image.ReTag(req.GetRetag().GetFrom(), req.GetRetag().GetTo())
		return image
	})
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	if req.Pull {
		buf.Write([]byte(images.GeneratePullScript()))
	}

	if req.Load {
		buf.Write([]byte(images.GenerateLoadScript()))
	}

	if req.Save {
		buf.Write([]byte(images.GenerateSaveScript()))
	}

	resp = &proto.GenerateScriptResponse{
		Out: buf.String(),
	}

	return
}
