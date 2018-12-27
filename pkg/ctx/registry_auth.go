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
	ra = &RegistryAuth{}
	ra.Username = env.Lookup(captain.EnvRegistryAuthUsername, captain.DefaultRegistryAuthUsername)
	ra.Password = env.Lookup(captain.EnvRegistryAuthPassword, captain.DefaultRegistryAuthPassword)
	return
}


func (e *RegistryAuth) AddFlags(f *pflag.FlagSet) {
	return
}
