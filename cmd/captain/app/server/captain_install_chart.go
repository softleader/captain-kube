package server

import (
	"context"
	"github.com/softleader/captain-kube/pkg/caplet"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/proto"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (s *CaptainServer) InstallChart(c context.Context, req *proto.InstallChartRequest) (resp *proto.InstallChartResponse, err error) {
	tmp, err := ioutil.TempDir(os.TempDir(), "install-chart-")
	if err != nil {
		return nil, err
	}
	chartPath := filepath.Join(tmp, req.GetChart().GetFileName())
	if err := ioutil.WriteFile(chartPath, req.GetChart().GetContent(), 0644); err != nil {
		return nil, err
	}

	i, err := chart.NewInstaller(req.GetK8S(), chartPath)
	if err != nil {
		return nil, err
	}
	if err := i.Install(); err != nil {
		return nil, err
	}

	if req.GetSync() {
		endpoints, err := s.lookupCaplets()
		if err != nil {
			return nil, err
		}
		tpls, err := chart.LoadArchive(s.out, chartPath)
		if err != nil {
			return nil, err
		}
		caplet.PullImage(s.out, endpoints, newPullImageRequest(tpls, req.GetRegistryAuth()), req.GetTimeout())
	}

	return
}

func newPullImageRequest(tpls chart.Templates, auth *proto.RegistryAuth) (req *proto.PullImageRequest) {
	req = &proto.PullImageRequest{
		RegistryAuth: auth,
	}
	for _, tpl := range tpls {
		for _, img := range tpl {
			req.Images = append(req.Images, &proto.Image{
				Host: img.Host,
				Repo: img.Repo,
				Tag:  img.Tag,
			})
		}
	}
	return
}
