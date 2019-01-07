package app

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/ctx"
	"github.com/softleader/captain-kube/pkg/dockerd"
	"github.com/spf13/cobra"
)

type pruneCmd struct {
	endpoint *ctx.Endpoint
}

func newPruneCmd() *cobra.Command {
	c := pruneCmd{
		endpoint: activeCtx.Endpoint,
	}

	cmd := &cobra.Command{
		Use:   "prune",
		Short: "docker system prune to all node",
		Long:  "docker system prune to all node",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.run()
		},
	}

	f := cmd.Flags()
	c.endpoint.AddFlags(f)

	return cmd
}

func (c *pruneCmd) run() error {
	if err := dockerd.Prune(logrus.StandardLogger()); err != nil {
		return err
	}
	if err := c.endpoint.Validate(); err == nil {
		return captain.Prune(logrus.StandardLogger(), c.endpoint.String(), settings.Verbose, settings.Color, settings.Timeout)
	}
	return nil
}
