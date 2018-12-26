package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/caplet"
	"github.com/softleader/captain-kube/pkg/color"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/version"
	"net"
)

var ErrNonCapletDaemonFound = fmt.Errorf("non caplet daemon found")

type CaptainServer struct {
	Metadata  *version.BuildMetadata
	Hostname  string
	Endpoints []string
	Port      int
	K8s       string
}

func (s *CaptainServer) lookupCaplet(log *logrus.Logger, colored bool) (endpoints caplet.Endpoints, err error) {
	if len(s.Endpoints) > 0 {
		log.Debugf("server has specified endpoint(s) for %q, skip dynamically lookup", s.Endpoints)
	} else {
		log.Debugf("looking up %q", s.Hostname)
		if s.Endpoints, err = net.LookupHost(s.Hostname); err != nil {
			return
		}
	}
	if len(s.Endpoints) == 0 {
		return nil, ErrNonCapletDaemonFound
	}
	for _, ep := range s.Endpoints {
		endpoints = append(endpoints, caplet.NewEndpoint(ep, s.Port))
	}
	if colored {
		for i, color := range color.Pick(len(endpoints)) {
			endpoints[i].Color = color
		}
	}
	log.Debugf("found %v caplet(s) daemon:", len(endpoints))
	if log.IsLevelEnabled(logrus.DebugLevel) {
		for _, e := range endpoints {
			log.Debugln(e.String())
		}
	}
	return
}

func newPullImageRequest(tpls chart.Templates, retag *proto.ReTag, auth *proto.RegistryAuth) (req *proto.PullImageRequest) {
	req = &proto.PullImageRequest{
		RegistryAuth: auth,
	}
	for _, tpl := range tpls {
		for _, img := range tpl {
			i := &proto.Image{
				Host: img.Host,
				Repo: img.Repo,
				Tag:  img.Tag,
			}
			if i.Host == retag.From {
				i.Host = retag.To
			}
			req.Images = append(req.Images, i)
		}
	}
	return
}
