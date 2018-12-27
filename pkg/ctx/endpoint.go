package ctx

import (
	"errors"
	"fmt"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/spf13/pflag"
)

type Endpoint struct {
	Host string
	Port int
}

func newEndpointFromEnv() (e *Endpoint) {
	e = &Endpoint{}
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

func (e *Endpoint) AddFlags(f *pflag.FlagSet) {
	f.StringVarP(&e.Host, "endpoint", "e", "", "specify the host of captain endpoint")
	f.IntVar(&e.Port, "endpoint-port", captain.DefaultPort, "specify the port of captain endpoint")
	return
}