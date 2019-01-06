package app

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/ctx"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/spf13/cobra"
	"strings"
)

type deleteCmd struct {
	chartName, chartVersion string
	helmTiller              *ctx.HelmTiller // helm tiller
	endpoint                *ctx.Endpoint   // captain 的 endpoint ip
}

func newDeleteCmd() *cobra.Command {
	c := deleteCmd{
		endpoint:   activeCtx.Endpoint,
		helmTiller: activeCtx.HelmTiller,
	}

	cmd := &cobra.Command{
		Use:   "delete CHART_NAME",
		Short: "delete helm-chart",
		Long:  "delete helm chart",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c.chartName = args[0]
			// do some validation check
			if err := c.endpoint.Validate(); err != nil {
				return err
			}
			// apply some default value
			if te := strings.TrimSpace(c.helmTiller.Endpoint); len(te) == 0 {
				c.helmTiller.Endpoint = c.endpoint.Host
			}
			return c.run()
		},
	}

	f := cmd.Flags()
	f.StringVarP(&c.chartVersion, "chart-version", "V", c.chartVersion, "the version of the helm chart to delete")
	c.endpoint.AddFlags(f)
	c.helmTiller.AddFlags(f)

	return cmd
}

func (c *deleteCmd) run() error {
	request := proto.DeleteChartRequest{
		Timeout: settings.timeout,
		Verbose: settings.verbose,
		Tiller: &proto.Tiller{
			Endpoint:          c.helmTiller.Endpoint,
			Username:          c.helmTiller.Username,
			Password:          c.helmTiller.Password,
			Account:           c.helmTiller.Account,
			SkipSslValidation: c.helmTiller.SkipSslValidation,
		},
		ChartName:    c.chartName,
		ChartVersion: c.chartVersion,
	}
	if err := captain.DeleteChart(logrus.StandardLogger(), c.endpoint.String(), &request, settings.timeout); err != nil {
		return err
	}
	return nil
}
