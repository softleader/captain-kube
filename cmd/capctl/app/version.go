package app

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	endpoint := activeCtx.Endpoint
	var full bool
	cmd := &cobra.Command{
		Use:   "version",
		Short: "印出 capctl, captain 及 caplet 的版本",
		Long:  "print capctl, captain, and caplet version",
		RunE: func(cmd *cobra.Command, args []string) error {
			logrus.Infoln(metadata.String())
			if endpoint.Specified() {
				if err := captain.Version(logrus.StandardLogger(), endpoint.String(), full, settings.Color, settings.TimeoutDuration()); err != nil {
					return err
				}
			}
			return nil
		},
	}

	f := cmd.Flags()
	f.BoolVar(&full, "full", false, "print full version number and commit hash")
	endpoint.AddFlags(f)

	return cmd
}
