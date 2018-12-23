package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/softleader/captain-kube/cmd/capctl/install"
	"github.com/softleader/captain-kube/cmd/capctl/script"
	"github.com/softleader/captain-kube/pkg/utils"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	var log *logrus.Logger
	var verbose bool

	cmd := &cobra.Command{
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

	f := cmd.Flags()
	f.BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")

	cmd.AddCommand(
		install.NewCmd(log),
		script.NewCmd(log),
	)

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}

}
