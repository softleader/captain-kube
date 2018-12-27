package ctx

import (
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/env"
	"github.com/spf13/pflag"
)

type Tiller struct {
	Endpoint          string // helm tiller 的 ip, 若沒改預設為 endpoint
	Username          string // helm tiller 的使用者
	Password          string // helm tiller 的密碼
	Account           string // helm tiller 的帳號
	SkipSslValidation bool
}

func newTillerFromEnv() (t *Tiller) {
	t.Endpoint = env.Lookup(captain.EnvTillerEndpoint, "")
	t.Username = env.Lookup(captain.EnvTillerUsername, captain.DefaultTillerUsername)
	t.Password = env.Lookup(captain.EnvTillerPassword, captain.DefaultTillerPassword)
	t.Account = env.Lookup(captain.EnvTillerAccount, captain.DefaultTillerAccount)
	t.SkipSslValidation = env.LookupBool(captain.EnvTillerSkipSslValidation, captain.DefaultTillerSkipSslValidation)
	return
}

func (e *Tiller) AddFlags(f *pflag.FlagSet) {
	return
}
