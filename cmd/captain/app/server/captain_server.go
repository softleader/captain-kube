package server

import (
	"fmt"
	"github.com/softleader/captain-kube/pkg/caplet"
	"io"
	"net"
)

var ErrNonCapletDaemonFound = fmt.Errorf("non caplet daemon found")

type CaptainServer struct {
	out       io.Writer
	hostname  string
	endpoints []string
	port      int
}

func NewCaptainServer(out io.Writer, hostname string, endpoints []string, port int) *CaptainServer {
	return &CaptainServer{
		out:       out,
		hostname:  hostname,
		endpoints: endpoints,
		port:      port,
	}
}

func (s *CaptainServer) retrieveCaplets() (endpoints []*caplet.Endpoint, err error) {
	if len(s.endpoints) == 0 {
		if s.endpoints, err = net.LookupHost(s.hostname); err != nil {
			return
		}
	}
	if len(s.endpoints) == 0 {
		return nil, ErrNonCapletDaemonFound
	}
	for _, ep := range s.endpoints {
		endpoints = append(endpoints, &caplet.Endpoint{
			Target: ep,
			Port:   s.port,
		})
	}
	return
}
