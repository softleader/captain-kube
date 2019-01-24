package chart

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/kubectl"
	"github.com/softleader/captain-kube/pkg/proto"
)

type Deleter interface {
	Delete(log *logrus.Logger) error
}

func NewDeleter(k8s *kubectl.KubeVersion, tiller *captainkube_v2.Tiller, chartName, chartVersion string) (Deleter, error) {
	if k8s.ServerVersion.IsICP() {
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

	if k8s.ServerVersion.IsGCP() {
		return &gcpDeleter{
			endpoint:     tiller.GetEndpoint(),
			chartName:    chartName,
			chartVersion: chartVersion,
		}, nil
	}

	return nil, fmt.Errorf("unsupported kubernetes vendor: %v", k8s.ServerVersion.GitCommit)
}
