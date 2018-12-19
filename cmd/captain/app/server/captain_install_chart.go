package server

import (
	"context"
	"github.com/softleader/captain-kube/pkg/caplet"
	"github.com/softleader/captain-kube/pkg/proto"
)

func (s *CaptainServer) InstallChart(c context.Context, req *proto.InstallChartRequest) (*proto.InstallChartResponse, error) {

	// TODO:
	// 1. 解開壓縮檔
	// 2. 執行 helm template, 取得要部署的檔案
	// 3. parse 出所有 image
	// 4. parallel call caplet pull

	// TODO: 從 chart 來
	var images []*proto.Image
	images = append(images, &proto.Image{
		Host: "",
		Repo: "alpine",
		Tag:  "",
	})

	caplet.PullImage(s.out, s.endpoints, s.port, &proto.PullImageRequest{
		Images: images,
	}, req.GetTimeout())

	return nil, nil
}
