package app

import (
	"github.com/Sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/utils"
	"github.com/softleader/captain-kube/pkg/version"
	"github.com/spf13/cobra"
)

func NewRootCmd(args []string, metadata *version.BuildMetadata) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "capctl",
		Short:        "the Captain-Kube command line interface",
		Long:         "The command line interface against Captain-Kube services",
		SilenceUsage: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logrus.SetFormatter(&utils.PlainFormatter{})
			if settings.verbose {
				logrus.SetLevel(logrus.DebugLevel)
			}
		},
	}

	flags := cmd.PersistentFlags()
	addGlobalFlags(flags)

	cmd.AddCommand(
		newInstallCmd(),
		newScriptCmd(),
		newPruneCmd(),
		newVersionCmd(metadata),
	)

	flags.Parse(args)

	return cmd
}
