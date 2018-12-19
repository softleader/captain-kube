package server

import (
	"context"
	"github.com/softleader/captain-kube/pkg/proto"
)

func (s *CaptainServer) InstallChart(context.Context, *proto.InstallChartRequest) (*proto.InstallChartResponse, error) {

	// TODO:
	// 1. 解開壓縮檔
	// 2. 執行 helm template, 取得要部署的檔案
	// 3. parse 出所有 image
	// 4. parallel call caplet pull

	return nil, nil
}
