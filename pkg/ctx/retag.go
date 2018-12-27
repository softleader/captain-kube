package ctx

import (
	"github.com/softleader/captain-kube/pkg/capui"
	"github.com/softleader/captain-kube/pkg/env"
	"github.com/spf13/pflag"
)

type ReTag struct {
	from string
	to   string
}

// 將系統 env 載入
func (rt *ReTag) ExpandEnv() {
	rt.from = env.Lookup(capui.EnvReTagFrom, capui.DefaultReTagFrom)
	rt.to = env.Lookup(capui.EnvReTagTo, capui.DefaultReTagTo)
}

func (rt *ReTag) AddFlags(f *pflag.FlagSet) {
	f.StringVarP(&rt.from, "retag-from", "f", rt.from, "specify re-tagging image host from")
	f.StringVarP(&rt.to, "retag-to", "t", rt.to, "specify re-tagging image host to")
	return
}
