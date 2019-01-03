package app

import (
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/ctx"
	"github.com/softleader/captain-kube/pkg/utils"
	"github.com/softleader/captain-kube/pkg/version"
	"github.com/spf13/cobra"
)

func NewRootCmd(args []string, metadata *version.BuildMetadata) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:          "capctl",
		Short:        "the captain-kube command line interface",
		Long:         "The command line interface against Captain-Kube services",
		SilenceUsage: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logrus.SetFormatter(&utils.PlainFormatter{})
			logrus.SetOutput(colorable.NewColorableStdout()) // for windows color output
			if settings.verbose {
				logrus.SetLevel(logrus.DebugLevel)
			}
		},
	}

	ctxs, err := ctx.LoadContextsFromEnv(logrus.StandardLogger())
	if err != nil {
		if err != ctx.ErrMountVolumeNotExist {
			return nil, err
		}
		ctxs = ctx.PlainContexts
	}
	activeCtx, err := ctxs.GetActiveExpandEnv()
	if err != nil {
		if err != ctx.ErrNoActiveContextPresent {
			return nil, err
		}
		activeCtx = ctx.NewContextFromEnv()
	}

	flags := cmd.PersistentFlags()
	addGlobalFlags(flags)

	cmd.AddCommand(
		newInstallCmd(activeCtx),
		newDeleteCmd(activeCtx),
		newScriptCmd(activeCtx),
		newPruneCmd(activeCtx),
		newVersionCmd(activeCtx, metadata),
		newCtxCmd(ctxs),
	)

	flags.Parse(args)

	return cmd, nil
}
