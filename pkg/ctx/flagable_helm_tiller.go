package ctx

import (
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/env"
	"github.com/spf13/pflag"
)

// HelmTiller 封裝了 helm tiller 相關資訊
type HelmTiller struct {
	Endpoint          string // helm tiller 的 ip, 若沒改預設為 endpoint
	Username          string // helm tiller 的使用者
	Password          string // helm tiller 的密碼
	Account           string // helm tiller 的帳號
	SkipSslValidation bool
}

func newHelmTillerFromEnv() (ht *HelmTiller) {
	ht = &HelmTiller{
		Endpoint:          env.Lookup(captain.EnvTillerEndpoint, ""),
		Username:          env.Lookup(captain.EnvTillerUsername, captain.DefaultTillerUsername),
		Password:          env.Lookup(captain.EnvTillerPassword, captain.DefaultTillerPassword),
		Account:           env.Lookup(captain.EnvTillerAccount, captain.DefaultTillerAccount),
		SkipSslValidation: env.LookupBool(captain.EnvTillerSkipSslValidation, captain.DefaultTillerSkipSslValidation),
	}
	return
}

// AddFlags 加入 flags
func (ht *HelmTiller) AddFlags(f *pflag.FlagSet) {
	f.StringVar(&ht.Endpoint, "tiller", ht.Endpoint, "specify the endpoint of helm tiller, override $"+captain.EnvTillerEndpoint)
	f.StringVar(&ht.Username, "tiller-user", ht.Username, "specify the username of helm tiller, override $"+captain.EnvTillerUsername)
	f.StringVar(&ht.Password, "tiller-pass", ht.Password, "specify the password of helm tiller, override $"+captain.EnvTillerPassword)
	f.StringVar(&ht.Account, "tiller-account", ht.Account, "specify the account of helm tiller, override $"+captain.EnvTillerAccount)
	f.BoolVar(&ht.SkipSslValidation, "tiller-skip-ssl", ht.SkipSslValidation, "skip the ssl validation of helm tiller, override $"+captain.EnvTillerSkipSslValidation)
	return
}
