package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/softleader/captain-kube/cmd/capctl/cmd"
	"github.com/softleader/captain-kube/pkg/utils"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

func main() {
	var log *logrus.Logger
	verbose, _ := strconv.ParseBool(os.Getenv("SL_VERBOSE"))

	rootCmd := &cobra.Command{
		Use:          "cap",
		Short:        "captain cli",
		Long:         "command intrface for captain",
		SilenceUsage: true,
		Run: func(cmd *cobra.Command, args []string) {
			log = logrus.New()
			log.SetOutput(cmd.OutOrStdout())
			log.SetFormatter(&utils.PlainFormatter{})
			if verbose {
				log.SetLevel(logrus.DebugLevel)
			}
		},
	}

	pf := rootCmd.PersistentFlags()
	pf.BoolVarP(&verbose, "verbose", "v", verbose, "enable verbose output")

	rootCmd.AddCommand(
		cmd.NewInstallCmd(log, verbose),
		cmd.NewScriptCmd(log, verbose),
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

}
