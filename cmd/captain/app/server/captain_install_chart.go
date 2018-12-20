package server

import (
	"context"
	"github.com/softleader/captain-kube/pkg/caplet"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/proto"
)

func (s *CaptainServer) InstallChart(c context.Context, req *proto.InstallChartRequest) (resp *proto.InstallChartResponse, err error) {
	chartFile := req.GetChart().GetFileName() // TODO

	if req.GetSync() {
		if tpls, err := chart.LoadArchive(s.out, chartFile); err != nil {
			return nil, err
		} else {
			caplet.PullImage(s.out, s.endpoints, s.port, newPullImageRequest(tpls), req.GetTimeout())
		}
	}

	// TODO: how to get caplet out?
	resp = &proto.InstallChartResponse{
		Out: "looks good!!?",
	}
	return
}

func newPullImageRequest(tpls chart.Templates) (req *proto.PullImageRequest) {
	req = &proto.PullImageRequest{}
	for _, imgs := range tpls {
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
