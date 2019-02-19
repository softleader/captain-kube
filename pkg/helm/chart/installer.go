package chart

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/kubectl"
	"github.com/softleader/captain-kube/pkg/proto"
)

type Installer interface {
	Install(log *logrus.Logger) error
}

func NewInstaller(log *logrus.Logger, k8s *kubectl.KubeVersion, tiller *captainkube_v2.Tiller, chart string) (Installer, error) {
	if tiller.GetEndpoint() == "" {
		return nil, fmt.Errorf("tiller endpoint is required")
	}

	if k8s.Server.IsICP() {
		log.Debugf("creating ICP chart installer")
		return &icpInstaller{
			endpoint:          tiller.GetEndpoint(),
			chart:             chart,
			username:          tiller.GetUsername(),
			password:          tiller.GetPassword(),
			account:           tiller.GetAccount(),
			skipSslValidation: tiller.GetSkipSslValidation(),
		}, nil
	}

	if k8s.Server.IsGCP() {
		log.Debugf("creating GCP chart installer")
		return &gcpInstaller{
			endpoint: tiller.GetEndpoint(),
			chart:    chart,
		}, nil
	}

	return nil, fmt.Errorf("unsupported kubernetes vendor: %v", k8s.Server.GitVersion)
}