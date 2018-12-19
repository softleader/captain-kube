package server

import (
	"context"
	"github.com/softleader/captain-kube/pkg/arc"
	"github.com/softleader/captain-kube/pkg/caplet"
	"github.com/softleader/captain-kube/pkg/helm"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/proto"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (s *CaptainServer) InstallChart(c context.Context, req *proto.InstallChartRequest) (resp *proto.InstallChartResponse, err error) {
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

	caplet.PullImage(s.out, s.endpoints, s.port, newPullImageRequest(images), req.GetTimeout())

	// TODO: how to get caplet out?
	resp = &proto.InstallChartResponse{
		Out: "looks good!!?",
	}
	return
}

func newPullImageRequest(images chart.Images) (req *proto.PullImageRequest) {
	req = &proto.PullImageRequest{}
	for _, imgs := range images {
		for _, img := range imgs {
			req.Images = append(req.Images, &proto.Image{
				Host: img.Host,
				Repo: img.Repo,
				Tag:  img.Tag,
			})
		}
	}
	return
}
