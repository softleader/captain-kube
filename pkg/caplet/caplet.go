package caplet

import (
	"fmt"
	"github.com/softleader/captain-kube/pkg/color"
	"sync"
)

const (
	EnvPort         = "CAPLET_PORT"
	EnvHostname     = "CAPLET_HOSTNAME"
	DefaultPort     = 50051
	DefaultHostname = "caplet"
)

type Endpoint struct {
	Target string
	Port   int
	Color  func([]byte) []byte // output 塗色
}

func NewEndpoint(target string, port int) *Endpoint {
	return &Endpoint{
		Target: target,
		Port:   port,
		Color:  color.Plain,
	}
}

func (e *Endpoint) String() string {
	return fmt.Sprintf("%s:%v", e.Target, e.Port)
}

type Endpoints []*Endpoint

func (endpoints Endpoints) Each(consumer func(e *Endpoint)) {
	var wg sync.WaitGroup
	for _, ep := range endpoints {
		wg.Add(1)
		go func(ep *Endpoint) {
			defer wg.Done()
			consumer(ep)
		}(ep)
	}
	wg.Wait()
}
