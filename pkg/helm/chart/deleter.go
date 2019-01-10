package chart

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/proto"
)

type Deleter interface {
	Delete(log *logrus.Logger) error
}

func NewDeleter(k8s string, tiller *tw_com_softleader_captainkube.Tiller, chartName, chartVersion string) (Deleter, error) {
	switch k8s {
	case "icp":
		return &icpDeleter{
			endpoint:          tiller.GetEndpoint(),
			username:          tiller.GetUsername(),
			password:          tiller.GetPassword(),
			account:           tiller.GetAccount(),
			skipSslValidation: tiller.GetSkipSslValidation(),
			chartName:         chartName,
			chartVersion:      chartVersion,
		}, nil
	case "gcp":
		return &gcpDeleter{
			endpoint:     tiller.GetEndpoint(),
			chartName:    chartName,
			chartVersion: chartVersion,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported kubernetes vendor: %v", k8s)
	}
}
