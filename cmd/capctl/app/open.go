package app

import (
	"github.com/pkg/browser"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/ctx"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/spf13/cobra"
)

type openCmd struct {
	endpoint *ctx.Endpoint
}

func newOpenCmd() *cobra.Command {
	c := openCmd{
		endpoint: activeCtx.Endpoint,
	}

	cmd := &cobra.Command{
		Use:   "open",
		Short: "以瀏覽器開啟 Kubernetes 管理介面",
		Long:  "Open endpoint Kubernetes console in browser",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.endpoint.Validate(); err != nil {
				return err
			}
			return c.run()
		},
	}

	f := cmd.Flags()
	c.endpoint.AddFlags(f)

	return cmd
}

func (c *openCmd) run() error {
	resp, err := captain.ConsoleURL(
		logrus.StandardLogger(),
		c.endpoint.String(),
		&captainkube_v2.ConsoleURLRequest{
			Host: c.endpoint.Host,
		},
		settings.Timeout)
	if err != nil {
		return err
	}
	url := resp.GetUrl()
	logrus.Debugf("opening %q in browser", url)
	return browser.OpenURL(url)
}
