package server

import (
	"fmt"
	"github.com/softleader/captain-kube/pkg/caplet"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/ver"
	"net"
)

var ErrNonCapletDaemonFound = fmt.Errorf("non caplet daemon found")

type CaptainServer struct {
	Metadata  *ver.BuildMetadata
	Hostname  string
	Endpoints []string
	Port      int
	K8s       string
}

func (s *CaptainServer) lookupCaplets() (endpoints caplet.Endpoints, err error) {
	if len(s.Endpoints) == 0 {
		if s.Endpoints, err = net.LookupHost(s.Hostname); err != nil {
			return
		}
	}
	if len(s.Endpoints) == 0 {
		return nil, ErrNonCapletDaemonFound
	}
	for _, ep := range s.Endpoints {
		endpoints = append(endpoints, &caplet.Endpoint{
			Target: ep,
			Port:   s.Port,
		})
	}
	return
}

func newPullImageRequest(tpls chart.Templates, auth *proto.RegistryAuth) (req *proto.PullImageRequest) {
	req = &proto.PullImageRequest{
		RegistryAuth: auth,
	}
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
