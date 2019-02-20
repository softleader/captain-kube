package app

import (
	"bytes"
	"fmt"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/ctx"
	"github.com/softleader/captain-kube/pkg/release"
	"github.com/softleader/captain-kube/pkg/utils"
	"github.com/spf13/cobra"
	"os"
	"text/template"
)

var (
	ctxs      *ctx.Contexts
	activeCtx *ctx.Context
	metadata  *release.Metadata
	name      = "capctl"
)

// NewRootCmd 建立 root command
func NewRootCmd(args []string, m *release.Metadata) (*cobra.Command, error) {
	if err := initContext(); err != nil {
		return nil, err
	}
	metadata = m
	if cli, found := os.LookupEnv("SL_CLI"); found {
		name = fmt.Sprintf("%s cap", cli)
	}

	cmd := &cobra.Command{
		Use:          name,
		Short:        "the captain-kube command line interface",
		Long:         "The command line interface against Captain-Kube services",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if settings.Offline {
				return fmt.Errorf("can not run the command in offline mode")
			}
			logrus.SetFormatter(&utils.PlainFormatter{})
			logrus.SetOutput(colorable.NewColorableStdout()) // for windows color output
			if settings.Verbose {
				logrus.SetLevel(logrus.DebugLevel)
			}
			return nil
		},
	}

	flags := cmd.PersistentFlags()
	initGlobalFlags(activeCtx.Global, flags)

	cmd.AddCommand(
		newInstallCmd(),
		newDeleteCmd(),
		newScriptCmd(),
		newPruneCmd(),
		newVersionCmd(),
		newCtxCmd(),
		newPullCmd(),
		newReTagCmd(),
		newSyncCmd(),
		newSaveCmd(),
		newOpenCmd(),
		newRmiCmd(),
		newRmcCmd(),
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

func usage(tpl string) string {
	var buf bytes.Buffer
	parsed := template.Must(template.New("").Parse(tpl))
	err := parsed.Execute(&buf, name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return buf.String()
}
