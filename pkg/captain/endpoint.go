package captain

import (
	"errors"
	"fmt"
	"github.com/spf13/pflag"
)

type Endpoint struct {
	host string
	port int
}

func (e *Endpoint) String() string {
	return fmt.Sprintf("%s:%v", e.host, e.port)
}

func (e *Endpoint) specified() bool {
	return len(e.host) > 0 && e.port > 0
}

func (e *Endpoint) validate() error {
	if !e.specified() {
		return errors.New("endpoint is required")
	}
	return nil
}

func AddEndpointFlags(f *pflag.FlagSet) (e *Endpoint) {
	e = &Endpoint{}
	f.StringVarP(&e.host, "endpoint", "e", "", "specify the host of captain endpoint")
	f.IntVar(&e.port, "endpoint-port", DefaultPort, "specify the port of captain endpoint")
	return
}
