package app

import (
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/ctx"
	"github.com/softleader/captain-kube/pkg/utils"
	"github.com/softleader/captain-kube/pkg/version"
	"github.com/spf13/cobra"
)

var (
	ctxs      *ctx.Contexts
	activeCtx *ctx.Context
	metadata  *version.BuildMetadata
)

func NewRootCmd(args []string, m *version.BuildMetadata) (*cobra.Command, error) {
	if err := initContext(); err != nil {
		return nil, err
	}
	metadata = m

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

	flags := cmd.PersistentFlags()
	addGlobalFlags(flags)

	cmd.AddCommand(
		newInstallCmd(),
		newDeleteCmd(),
		newScriptCmd(),
		newPruneCmd(),
		newVersionCmd(),
		newCtxCmd(),
	)

	flags.Parse(args)

	return cmd, nil
}

func initContext() (err error) {
	if ctxs, err = ctx.LoadContextsFromEnv(logrus.StandardLogger()); err != nil {
		if err != ctx.ErrMountVolumeNotExist {
			return err
		}
		ctxs = ctx.PlainContexts
	}
	if activeCtx, err = ctxs.GetActiveExpandEnv(); err != nil {
		if err != ctx.ErrNoActiveContextPresent {
			return err
		}
		activeCtx = ctx.NewContextFromEnv()
	}
	return nil
}
