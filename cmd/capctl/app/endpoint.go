package app

import (
	"errors"
	"fmt"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/spf13/pflag"
	"strings"
)

type endpoint struct {
	host string
	port int
}

func (e *endpoint) String() string {
	return fmt.Sprintf("%s:%v", e.host, e.port)
}

func (e *endpoint) validate() error {
	if e := strings.TrimSpace(e.host); len(e) == 0 {
		return errors.New("endpoint is required")
	}
	return nil
}

func addEndpointFlags(f *pflag.FlagSet) (e *endpoint) {
	e = &endpoint{}
	f.StringVarP(&e.host, "endpoint", "e", "", "specify the host of captain endpoint")
	f.IntVar(&e.port, "endpoint-port", captain.DefaultPort, "specify the port of captain endpoint")
	return
}
