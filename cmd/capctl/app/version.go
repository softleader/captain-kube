package app

import (
	"github.com/Sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/version"
	"github.com/spf13/cobra"
)

func newVersionCmd(metadata *version.BuildMetadata) *cobra.Command {
	var endpoint *endpoint
	var short bool
	cmd := &cobra.Command{
		Use:   "version",
		Short: "print capctl, captain, and caplet version",
		Long:  "print capctl, captain, and caplet version",
		RunE: func(cmd *cobra.Command, args []string) error {
			logrus.Infoln(metadata.String(short))
			if endpoint.specified() {
				if err := captain.Version(logrus.StandardLogger(), endpoint.String(), short, settings.timeout); err != nil {
					return err
				}
			}
			return nil
		},
	}

	f := cmd.Flags()
	f.BoolVar(&short, "short", false, "print only the version number plus first 7 digits of the commit hash")
	endpoint = addEndpointFlags(f)

	return cmd
}
