package app

import (
	"github.com/softleader/captain-kube/pkg/dur"
	"github.com/spf13/pflag"
	"os"
	"strconv"
)

var settings = new(envSettings)

type envSettings struct {
	verbose bool
	timeout int64
}

func addGlobalFlags(fs *pflag.FlagSet) {
	settings.verbose, _ = strconv.ParseBool(os.Getenv("SL_VERBOSE"))
	fs.BoolVarP(&settings.verbose, "verbose", "v", settings.verbose, "enable verbose output")
	fs.Int64Var(&settings.timeout, "timeout", dur.DefaultDeadlineSecond, "timeout second communicating to captain services")
}