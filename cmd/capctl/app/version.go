package app

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/ctx"
	"github.com/softleader/captain-kube/pkg/version"
	"github.com/spf13/cobra"
)

func newVersionCmd(activeCtx *ctx.Context, metadata *version.BuildMetadata) *cobra.Command {
	endpoint := activeCtx.Endpoint
	var short bool
	cmd := &cobra.Command{
		Use:   "version",
		Short: "print capctl, captain, and caplet version",
		Long:  "print capctl, captain, and caplet version",
		RunE: func(cmd *cobra.Command, args []string) error {
			logrus.Infoln(metadata.String(short))
			if endpoint.Specified() {
				if err := captain.Version(logrus.StandardLogger(), endpoint.String(), short, settings.color, settings.timeout); err != nil {
					return err
				}
			}
			return nil
		},
	}

	f := cmd.Flags()
	f.BoolVarP(&short, "short", "s", false, "print only the version number plus first 7 digits of the commit hash")
	endpoint.AddFlags(f)

	return cmd
}
