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

func newReTagFromEnv() (rt *ReTag) {
	rt = &ReTag{
		From: env.Lookup(capui.EnvReTagFrom, capui.DefaultReTagFrom),
		To:   env.Lookup(capui.EnvReTagTo, capui.DefaultReTagTo),
	}
	return
}

func (rt *ReTag) AddFlags(f *pflag.FlagSet) {
	f.StringVarP(&rt.From, "retag-from", "f", rt.From, "specify re-tagging image host from, override $"+capui.EnvReTagFrom)
	f.StringVarP(&rt.To, "retag-to", "t", rt.To, "specify re-tagging image host to, override $"+capui.EnvReTagTo)
	return
}
