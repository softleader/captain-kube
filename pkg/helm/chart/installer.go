package chart

import (
	"fmt"
	"github.com/softleader/captain-kube/pkg/proto"
	"io"
)

type Installer interface {
	Install(out io.Writer) error
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
