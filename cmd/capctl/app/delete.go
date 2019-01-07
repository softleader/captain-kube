package app

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/ctx"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/spf13/cobra"
	"strings"
)

const (
	deleteHelp = `刪除 Helm Chart

使用 '--endpoint' 指定刪除的 Captain Endpoint

	$ {{.}} delete CHART_NAME -e CAPTAIN_ENDPOINT

若 Helm Tiller Server 不在 Captain-Kube 環境中, 可以傳入 '--tiller*' 開頭的 flag 設定 Tiller 相關資訊

	$ {{.}} delete CHART_NAME -e CAPTAIN_ENDPOINT --tiller TILLER_IP
	$ {{.}} delete CHART_NAME -e CAPTAIN_ENDPOINT --tiller TILLER_IP --tiller-skip-ssl=false

傳入 '--chart-version' 指定只刪除特定版本; 反之則刪除全部版本

	$ {{.}} delete CHART_NAME -e CAPTAIN_ENDPOINT -V 0.1.0
`
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
		Long:  usage(deleteHelp),
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
		Timeout: settings.Timeout,
		Verbose: settings.Verbose,
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
	if err := captain.DeleteChart(logrus.StandardLogger(), c.endpoint.String(), &request, settings.Timeout); err != nil {
		return err
	}
	return nil
}
