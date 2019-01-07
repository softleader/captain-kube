package app

import (
	"github.com/softleader/captain-kube/pkg/dur"
	"github.com/spf13/pflag"
	"os"
	"strconv"
)

var settings = new(envSettings)

type envSettings struct {
	offline bool
	verbose bool
	color   bool
	timeout int64
}

func addGlobalFlags(fs *pflag.FlagSet) {
	settings.offline, _ = strconv.ParseBool(os.Getenv("SL_OFFLINE"))
	settings.verbose, _ = strconv.ParseBool(os.Getenv("SL_VERBOSE"))
	fs.BoolVarP(&settings.offline, "offline", "o", settings.offline, "work offline")
	fs.BoolVarP(&settings.verbose, "verbose", "v", settings.verbose, "enable verbose output")
	fs.BoolVar(&settings.color, "color", settings.color, "colored caplet output")
	fs.Int64Var(&settings.timeout, "timeout", dur.DefaultDeadlineSecond, "timeout second communicating to captain services")
}
