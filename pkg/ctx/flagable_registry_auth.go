package ctx

import (
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/env"
	"github.com/spf13/pflag"
)

// RegistryAuth 封裝了 docker registry 授權資訊
type RegistryAuth struct {
	Username string // docker registry 的帳號
	Password string // docker registry 的密碼
}

func newRegistryAuthFromEnv() (ra *RegistryAuth) {
	ra = &RegistryAuth{
		Username: env.Lookup(captain.EnvRegistryAuthUsername, captain.DefaultRegistryAuthUsername),
		Password: env.Lookup(captain.EnvRegistryAuthPassword, captain.DefaultRegistryAuthPassword),
	}
	return
}

// AddFlags 加入 flags
func (ra *RegistryAuth) AddFlags(f *pflag.FlagSet) {
	f.StringVar(&ra.Username, "reg-user", ra.Username, "specify username of basic-auth for docker registry, override $"+captain.EnvRegistryAuthUsername)
	f.StringVar(&ra.Password, "reg-pass", ra.Password, "specify password of basic-auth for docker registry, override $"+captain.EnvRegistryAuthPassword)
	return
}
