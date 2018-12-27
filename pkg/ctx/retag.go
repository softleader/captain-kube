package ctx

import (
	"github.com/softleader/captain-kube/pkg/capui"
	"github.com/softleader/captain-kube/pkg/env"
	"github.com/spf13/pflag"
)

type ReTag struct {
	From string
	To   string
}

// 將系統 env 載入
func (rt *ReTag) ExpandEnv() {
	rt.From = env.Lookup(capui.EnvReTagFrom, capui.DefaultReTagFrom)
	rt.To = env.Lookup(capui.EnvReTagTo, capui.DefaultReTagTo)
}

func (rt *ReTag) AddFlags(f *pflag.FlagSet) {
	f.StringVarP(&rt.From, "retag-from", "f", rt.From, "specify re-tagging image host from")
	f.StringVarP(&rt.To, "retag-to", "t", rt.To, "specify re-tagging image host to")
	return
}
