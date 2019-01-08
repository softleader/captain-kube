package ctx

import (
	"github.com/softleader/captain-kube/pkg/dur"
	"github.com/spf13/pflag"
	"os"
	"strconv"
)

type Global struct {
	Offline bool
	Verbose bool
	Color   bool
	Timeout int64
}

func newGlobalFromEnv() (g *Global) {
	g = &Global{}
	g.Offline, _ = strconv.ParseBool(os.Getenv("SL_OFFLINE"))
	g.Verbose, _ = strconv.ParseBool(os.Getenv("SL_VERBOSE"))
	g.Timeout = dur.DefaultDeadlineSecond
	return
}

func (g *Global) AddFlags(f *pflag.FlagSet) {
	f.BoolVar(&g.Offline, "offline", g.Offline, "work offline")
	f.BoolVarP(&g.Verbose, "verbose", "v", g.Verbose, "enable verbose output")
	f.BoolVar(&g.Color, "color", g.Color, "colored caplet output")
	f.Int64Var(&g.Timeout, "timeout", dur.DefaultDeadlineSecond, "timeout second communicating to captain services")
	return
}
