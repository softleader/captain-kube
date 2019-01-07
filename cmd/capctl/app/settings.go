package app

import (
	"github.com/softleader/captain-kube/pkg/ctx"
	"github.com/spf13/pflag"
)

var settings *ctx.Global

func initGlobalFlags(g *ctx.Global, fs *pflag.FlagSet) {
	settings = g
	settings.AddFlags(fs)
}
