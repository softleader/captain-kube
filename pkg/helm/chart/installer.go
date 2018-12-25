package chart

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/proto"
)

type Installer interface {
	Install(log *logrus.Logger) error
}

func NewInstaller(k8s string, tiller *proto.Tiller, chart string) (Installer, error) {
	switch k8s {
	case "icp":
		return &icpInstaller{
			endpoint:          tiller.GetEndpoint(),
			chart:             chart,
			username:          tiller.GetUsername(),
			password:          tiller.GetPassword(),
			account:           tiller.GetAccount(),
			skipSslValidation: tiller.GetSkipSslValidation(),
		}, nil
	case "gcp":
		return &gcpInstaller{
			endpoint: tiller.GetEndpoint(),
			chart:    chart,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported kubernetes vendor: %v", k8s)
	}
}
