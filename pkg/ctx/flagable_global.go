package ctx

import (
	"github.com/softleader/captain-kube/pkg/dur"
	"github.com/spf13/pflag"
	"os"
	"strconv"
	"time"
)

type Global struct {
	Offline bool
	Verbose bool
	Color   bool
	Timeout string
}

func (g *Global) TimeoutDuration() time.Duration {
	return dur.Parse(g.Timeout)
}

func newGlobalFromEnv() (g *Global) {
	g = &Global{}
	g.Offline, _ = strconv.ParseBool(os.Getenv("SL_OFFLINE"))
	g.Verbose, _ = strconv.ParseBool(os.Getenv("SL_VERBOSE"))
	g.Timeout = dur.DefaultDeadline
	return
}

func (g *Global) AddFlags(f *pflag.FlagSet) {
	f.BoolVar(&g.Offline, "offline", g.Offline, "work offline")
	f.BoolVarP(&g.Verbose, "verbose", "v", g.Verbose, "enable verbose output")
	f.BoolVar(&g.Color, "color", g.Color, "colored caplet output")
	f.StringVar(&g.Timeout, "timeout", dur.DefaultDeadline, `timeout communicating to captain, supports units are "ms", "s", "m", "h"`)
	return
}
