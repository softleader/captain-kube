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

// 將系統 env 載入
func (e *Endpoint) ExpandEnv() {
	e.Host = env.Lookup(captain.EnvEndpoint, "")
	return
}

func (e *Endpoint) AddFlags(f *pflag.FlagSet) {
	f.StringVarP(&e.Host, "endpoint", "e", "", "specify the host of captain endpoint")
	f.IntVar(&e.Port, "endpoint-port", captain.DefaultPort, "specify the port of captain endpoint")
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
