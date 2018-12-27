package ctx

import (
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/env"
	"github.com/spf13/pflag"
)

type HelmTiller struct {
	Endpoint          string // helm tiller 的 ip, 若沒改預設為 endpoint
	Username          string // helm tiller 的使用者
	Password          string // helm tiller 的密碼
	Account           string // helm tiller 的帳號
	SkipSslValidation bool
}

// 將系統 env 載入
func (t *HelmTiller) ExpandEnv() {
	t.Endpoint = env.Lookup(captain.EnvTillerEndpoint, "")
	t.Username = env.Lookup(captain.EnvTillerUsername, captain.DefaultTillerUsername)
	t.Password = env.Lookup(captain.EnvTillerPassword, captain.DefaultTillerPassword)
	t.Account = env.Lookup(captain.EnvTillerAccount, captain.DefaultTillerAccount)
	t.SkipSslValidation = env.LookupBool(captain.EnvTillerSkipSslValidation, captain.DefaultTillerSkipSslValidation)
}

func (t *HelmTiller) AddFlags(f *pflag.FlagSet) {
	f.StringVar(&t.Endpoint, "tiller", t.Endpoint, "specify the endpoint of helm tiller")
	f.StringVar(&t.Username, "tiller-user", t.Username, "specify the username of helm tiller")
	f.StringVar(&t.Password, "tiller-pass", t.Password, "specify the password of helm tiller")
	f.StringVar(&t.Account, "tiller-account", t.Account, "specify the account of helm tiller")
	f.BoolVar(&t.SkipSslValidation, "tiller-skip-ssl", t.SkipSslValidation, "specify skip ssl validation of helm tiller")
	return
}
