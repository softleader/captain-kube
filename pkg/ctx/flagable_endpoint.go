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
	f.StringVarP(&e.Host, "endpoint", "e", e.Host, "specify the host of captain endpoint")
	f.IntVar(&e.Port, "endpoint-port", e.Port, "specify the port of captain endpoint")
	return
}

func (e *Endpoint) String() string {
	return fmt.Sprintf("%s:%v", e.Host, e.Port)
}

func (e *Endpoint) Specified() bool {
	return len(e.Host) > 0 && e.Port > 0
}

func (e *Endpoint) Validate() error {
	if !e.Specified() {
		return errors.New("endpoint is required")
	}
	return nil
}
