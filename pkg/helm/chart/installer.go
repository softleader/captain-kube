package chart

import (
	"fmt"
	"github.com/softleader/captain-kube/pkg/proto"
)

type Installer interface {
	Install() error
}

func NewInstaller(k8s *proto.K8S, chart string) (Installer, error) {
	switch v := k8s.GetVendor(); v {
	case proto.K8SVendor_Gcp:
		return &gcpInstaller{
			chart: chart,
		}, nil
	case proto.K8SVendor_Icp:
		return &icpInstaller{
			endpoint: k8s.GetEndpoint(),
			chart:    chart,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported K8s vendo: %s", v)
	}
}
