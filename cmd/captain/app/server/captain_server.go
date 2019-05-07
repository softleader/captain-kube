package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/caplet"
	"github.com/softleader/captain-kube/pkg/color"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/kubectl"
	pb "github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/release"
	"github.com/softleader/captain-kube/pkg/utils/tcp"
	"net"
)

// ErrNonCapletDaemonFound 代表沒有找到任何的 caplet daemons
var ErrNonCapletDaemonFound = fmt.Errorf("non caplet daemon found")

// CaptainServer 封裝了 captain server 的資訊
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
	for _, host := range hosts {
		// make a quick test to check the endpoint is available
		if !tcp.IsReachable(host, s.Port, 1) {
			log.Debugf("skipping unreachable caplet: %s:%v", host, s.Port)
			continue
		}
		endpoints = append(endpoints, caplet.NewEndpoint(host, s.Port))
	}
	if len(endpoints) == 0 {
		return nil, ErrNonCapletDaemonFound
	}
	if colored {
		for i, color := range color.Pick(len(endpoints)) {
			endpoints[i].Color = color
		}
	}
	log.Debugf("found %v caplet(s) daemon: %v", len(endpoints), endpoints)
	return
}

func newPullImageRequest(tpls chart.Templates, retag *pb.ReTag, auth *pb.RegistryAuth) (req *pb.PullImageRequest) {
	req = &pb.PullImageRequest{
		RegistryAuth: auth,
	}
	for _, tpl := range tpls {
		for _, img := range tpl {
			i := &pb.Image{
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
