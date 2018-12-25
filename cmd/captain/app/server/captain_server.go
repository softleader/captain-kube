package server

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/caplet"
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

func (s *CaptainServer) lookupInstances() (endpoints caplet.Endpoints, err error) {
	if len(s.Endpoints) > 0 {
		logrus.Debugf("server has specified endpoint %q, skip dynamically lookup", s.Endpoints)
	} else {
		logrus.Debugf("nslookup %q", s.Hostname)
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
	logrus.Debugf("found %v caplet(s) daemon:", len(endpoints))
	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		for _, e := range endpoints {
			logrus.Debugln(e.String())
		}
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
