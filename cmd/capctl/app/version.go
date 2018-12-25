package app

import (
	"github.com/Sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/ver"
	"github.com/spf13/cobra"
)

func newVersionCmd(metadata *ver.BuildMetadata) *cobra.Command {
	var endpoint *endpoint
	var short bool
	cmd := &cobra.Command{
		Use:   "version",
		Short: "print capctl, captain, and caplet version",
		Long:  "print capctl, captain, and caplet version",
		Run: func(cmd *cobra.Command, args []string) {
			logrus.Infoln(metadata.String(short))
			if endpoint.specified() {
				captain.Version(logrus.StandardLogger(), endpoint.String(), short, settings.timeout)
			}
		},
	}

	f := cmd.Flags()
	f.BoolVar(&short, "short", false, "print only the version number plus first 7 digits of the commit hash")
	endpoint = addEndpointFlags(f)

	return cmd
}
