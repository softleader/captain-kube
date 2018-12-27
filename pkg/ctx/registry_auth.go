package ctx

import (
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/env"
	"github.com/spf13/pflag"
)

type RegistryAuth struct {
	Username string // docker registry 的帳號
	Password string // docker registry 的密碼
}

// 將系統 env 載入
func (ra *RegistryAuth) ExpandEnv() {
	ra.Username = env.Lookup(captain.EnvRegistryAuthUsername, captain.DefaultRegistryAuthUsername)
	ra.Password = env.Lookup(captain.EnvRegistryAuthPassword, captain.DefaultRegistryAuthPassword)
}

func (ra *RegistryAuth) AddFlags(f *pflag.FlagSet) {
	f.StringVar(&ra.Username, "reg-user", ra.Username, "specify username of basic-auth for docker registry")
	f.StringVar(&ra.Password, "reg-pass", ra.Password, "specify password of basic-auth for docker registry")
	return
}
