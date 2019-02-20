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
		Short: "在每個 worker node 上刪除無用的 docker 資料",
		Long:  "docker system prune to current host and all worker nodes as well",
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
		return captain.CallPrune(logrus.StandardLogger(), c.endpoint.String(), settings.Verbose, settings.Color, settings.TimeoutDuration())
	}
	return nil
}
