package captain

import (
	"errors"
	"fmt"
	"github.com/spf13/pflag"
)

type Endpoint struct {
	Host string
	Port int
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

func AddEndpointFlags(f *pflag.FlagSet) (e *Endpoint) {
	e = &Endpoint{}
	f.StringVarP(&e.Host, "endpoint", "e", "", "specify the host of captain endpoint")
	f.IntVar(&e.Port, "endpoint-port", DefaultPort, "specify the port of captain endpoint")
	return
}
func AddEndpointFlagsWithDefault(f *pflag.FlagSet, e *Endpoint) *Endpoint {
	f.StringVarP(&e.Host, "endpoint", "e", e.Host, "specify the host of captain endpoint")
	f.IntVar(&e.Port, "endpoint-port", e.Port, "specify the port of captain endpoint")
	return e
}
