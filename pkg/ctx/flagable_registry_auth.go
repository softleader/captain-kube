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

func newRegistryAuthFromEnv() (ra *RegistryAuth) {
	ra = &RegistryAuth{
		Username: env.Lookup(captain.EnvRegistryAuthUsername, captain.DefaultRegistryAuthUsername),
		Password: env.Lookup(captain.EnvRegistryAuthPassword, captain.DefaultRegistryAuthPassword),
	}
	return
}

func (ra *RegistryAuth) AddFlags(f *pflag.FlagSet) {
	f.StringVar(&ra.Username, "reg-user", ra.Username, "specify username of basic-auth for docker registry")
	f.StringVar(&ra.Password, "reg-pass", ra.Password, "specify password of basic-auth for docker registry")
	return
}
