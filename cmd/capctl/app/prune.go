package app

import (
	"github.com/Sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/spf13/cobra"
)

type pruneCmd struct {
	endpoint *endpoint
}

func newPruneCmd() *cobra.Command {
	c := pruneCmd{}

	cmd := &cobra.Command{
		Use:   "prune",
		Short: "docker system prune to all node",
		Long:  "docker system prune to all node",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.endpoint.validate(); err != nil {
				return err
			}
			return c.run()
		},
	}

	f := cmd.Flags()
	c.endpoint = addEndpointFlags(f)

	return cmd
}

func (c *pruneCmd) run() error {
	return captain.Prune(logrus.StandardLogger(), c.endpoint.String(), &proto.PruneRequest{
		Verbose: settings.verbose,
		Timeout: settings.timeout,
	}, settings.timeout)
}
