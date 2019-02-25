package chart

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/kubectl"
	pb "github.com/softleader/captain-kube/pkg/proto"
)

// Deleter 定義了從 helm tiller 中刪除 chart 的介面
type Deleter interface {
	Delete(log *logrus.Logger) error
}

// NewDeleter 建立 Deleter
func NewDeleter(k8s *kubectl.KubeVersion, tiller *pb.Tiller, chartName, chartVersion string) (Deleter, error) {
	if k8s.Server.IsICP() {
		return &icpDeleter{
			endpoint:          tiller.GetEndpoint(),
			username:          tiller.GetUsername(),
			password:          tiller.GetPassword(),
			account:           tiller.GetAccount(),
			skipSslValidation: tiller.GetSkipSslValidation(),
			chartName:         chartName,
			chartVersion:      chartVersion,
		}, nil
	}

	if k8s.Server.IsGCP() {
		return &gcpDeleter{
			endpoint:     tiller.GetEndpoint(),
			chartName:    chartName,
			chartVersion: chartVersion,
		}, nil
	}

	return nil, fmt.Errorf("unsupported kubernetes vendor: %v", k8s.Server.GitVersion)
}
