package ctx

import (
	"errors"
	"fmt"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/env"
	"github.com/spf13/pflag"
)

// Endpoint 封裝了連線的基本資訊
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

// AddFlags 加入 flags
func (e *Endpoint) AddFlags(f *pflag.FlagSet) {
	f.StringVarP(&e.Host, "endpoint", "e", e.Host, "specify the host of captain endpoint, override $"+captain.EnvEndpoint)
	f.IntVar(&e.Port, "endpoint-port", e.Port, "specify the port of captain endpoint, override $"+captain.EnvPort)
	return
}

// String 回傳 Endpoint 資訊
func (e *Endpoint) String() string {
	return fmt.Sprintf("%s:%v", e.Host, e.Port)
}

// Specified 回傳 Endpoint 是否合法的被定義
func (e *Endpoint) Specified() bool {
	return e.Validate() == nil
}

// Validate 驗證 Endpoint 的值是否有問題
func (e *Endpoint) Validate() error {
	if len(e.Host) == 0 {
		return errors.New("endpoint is required")
	}
	if e.Port == 0 {
		return errors.New("endpoint port is required")
	}
	return nil
}
