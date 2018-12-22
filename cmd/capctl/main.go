package main

import (
	"github.com/softleader/captain-kube/cmd/capctl/install"
	"github.com/softleader/captain-kube/cmd/capctl/script"
	"github.com/softleader/captain-kube/pkg/logger"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	var log *logger.Logger
	var verbose bool

	cmd := &cobra.Command{
		Use:          "cap",
		Short:        "captain cli",
		Long:         "command intrface for captain",
		SilenceUsage: true,
		Run: func(cmd *cobra.Command, args []string) {
			log = logger.New(cmd.OutOrStdout()).WithVerbose(verbose)
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
