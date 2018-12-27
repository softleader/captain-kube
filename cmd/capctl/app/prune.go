package app

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/ctx"
	"github.com/spf13/cobra"
)

type pruneCmd struct {
	endpoint *ctx.Endpoint
}

func newPruneCmd(activeCtx *ctx.Context) *cobra.Command {
	c := pruneCmd{
		endpoint: activeCtx.Endpoint,
	}

	cmd := &cobra.Command{
		Use:   "prune",
		Short: "docker system prune to all node",
		Long:  "docker system prune to all node",
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

func (c *pruneCmd) run() error {
	return captain.Prune(logrus.StandardLogger(), c.endpoint.String(), settings.verbose, settings.color, settings.timeout)
}
