package app

import (
	"github.com/Sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/utils"
	"github.com/spf13/cobra"
)

func NewRootCmd(args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "capctl",
		Short:        "captain-kube cli",
		Long:         "The command line interface for captain-kube system",
		SilenceUsage: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logrus.SetOutput(cmd.OutOrStdout())
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
	)

	flags.Parse(args)

	return cmd
}
