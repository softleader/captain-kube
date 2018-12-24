package caplet

import (
	"fmt"
	"strings"
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
}

type Endpoints []*Endpoint

func (endpoints Endpoints) Each(consumer func(e *Endpoint) error) error {
	ch := make(chan error, len(endpoints))
	var wg sync.WaitGroup
	for _, ep := range endpoints {
		wg.Add(1)
		go func(ep *Endpoint) {
			defer wg.Done()
			ch <- consumer(ep)
		}(ep)
	}
	wg.Wait()
	close(ch)
	var errors []string
	for e := range ch {
		if e != nil {
			errors = append(errors, e.Error())
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "\n"))
	}
	return nil
}
