package ctx

import (
	"github.com/softleader/captain-kube/pkg/capui"
	"github.com/softleader/captain-kube/pkg/env"
	"github.com/spf13/pflag"
	"strings"
)

// HelmChart 封裝了跟 Helm Chart 有關的 command
type HelmChart struct {
	Set []string
}

func newHelmChartFromEnv() (c *HelmChart) {
	c = &HelmChart{
		Set: strings.Split(env.Lookup(capui.EnvChartSet, capui.DefaultChartSet), ","),
	}
	return
}

// AddFlags 加入 flags
func (c *HelmChart) AddFlags(f *pflag.FlagSet) {
	f.StringArrayVar(&c.Set, "set", c.Set, "set & overwrite values to rendered templates, will not be affected the original Chart file (.tgz) (can specify multiple or separate values with commas: key1=val1,key2=val2)")
	return
}
