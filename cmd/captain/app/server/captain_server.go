package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/caplet"
	"github.com/softleader/captain-kube/pkg/color"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/kubectl"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/release"
	"net"
)

var ErrNonCapletDaemonFound = fmt.Errorf("non caplet daemon found")

type CaptainServer struct {
	Metadata  *release.Metadata
	Hostname  string
	Endpoints []string
	Port      int
	K8s       *kubectl.KubeVersion
}

func (s *CaptainServer) lookupCaplet(log *logrus.Logger, colored bool) (endpoints caplet.Endpoints, err error) {
	var hosts []string
	if len(s.Endpoints) > 0 {
		log.Debugf("server has specified endpoint(s) for %q, skip dynamically lookup", s.Endpoints)
		hosts = s.Endpoints
	} else {
		log.Debugf("looking up %q", s.Hostname)
		if hosts, err = net.LookupHost(s.Hostname); err != nil {
			return
		}
	}
	if len(hosts) == 0 {
		return nil, ErrNonCapletDaemonFound
	}
	for _, host := range hosts {
		endpoints = append(endpoints, caplet.NewEndpoint(host, s.Port))
	}
	if colored {
		for i, color := range color.Pick(len(endpoints)) {
			endpoints[i].Color = color
		}
	}
	log.Debugf("found %v caplet(s) daemon: %v", len(endpoints), endpoints)
	return
}

func newPullImageRequest(tpls chart.Templates, retag *captainkube_v2.ReTag, auth *captainkube_v2.RegistryAuth) (req *captainkube_v2.PullImageRequest) {
	req = &captainkube_v2.PullImageRequest{
		RegistryAuth: auth,
	}
	for _, tpl := range tpls {
		for _, img := range tpl {
			i := &captainkube_v2.Image{
				Host: img.Host,
				Repo: img.Repo,
				Tag:  img.Tag,
			}
			if i.Host == retag.From && len(retag.To) > 0 {
				i.Host = retag.To
			}
			req.Images = append(req.Images, i)
		}
	}
	return
}
