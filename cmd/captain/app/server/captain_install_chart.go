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
		endpoints, err := s.retrieveCaplets()
		if err != nil {
			return nil, err
		}
		tpls, err := chart.LoadArchive(s.out, chartFile)
		if err != nil {
			return nil, err
		}
		caplet.PullImage(s.out, endpoints, newPullImageRequest(tpls), req.GetTimeout())
	}

	// TODO: how to get caplet out?
	resp = &proto.InstallChartResponse{
		Out: "looks good!!?",
	}
	return
}

func newPullImageRequest(tpls chart.Templates) (req *proto.PullImageRequest) {
	req = &proto.PullImageRequest{}
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
