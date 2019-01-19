package ctx

import (
	"errors"
	"fmt"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/env"
	"github.com/spf13/pflag"
)

type Endpoint struct {
	Host string
	Port int
}

func newEndpointFromEnv() (e *Endpoint) {
	e = &Endpoint{
		Host: env.Lookup(captain.EnvEndpoint, ""),
		Port: env.LookupInt(captain.EnvPort, captain.DefaultPort),
	}
	return
}

func (e *Endpoint) AddFlags(f *pflag.FlagSet) {
	f.StringVarP(&e.Host, "endpoint", "e", e.Host, "specify the host of captain endpoint, override $"+captain.EnvEndpoint)
	f.IntVar(&e.Port, "endpoint-port", e.Port, "specify the port of captain endpoint, override $"+captain.EnvPort)
	return
}

func (e *Endpoint) String() string {
	return fmt.Sprintf("%s:%v", e.Host, e.Port)
}

func (e *Endpoint) Specified() bool {
	return e.Validate() == nil
}

func (e *Endpoint) Validate() error {
	if len(e.Host) == 0 {
		return errors.New("endpoint is required")
	}
	if e.Port == 0 {
		return errors.New("endpoint port is required")
	}
	return nil
}
