package main

import (
	"github.com/softleader/captain-kube/cmd/capctl/install"
	"github.com/softleader/captain-kube/cmd/capctl/script"
	"github.com/spf13/cobra"
	"io"
	"os"
)

type rootCmd struct {
	out  io.Writer
	tags string
}

func main() {
	c := rootCmd{}

	cmd := &cobra.Command{
		Use:          "cap",
		Short:        "captain cli",
		Long:         "command intrface for captain",
		SilenceUsage: true,
	}

	c.out = cmd.OutOrStdout()

	cmd.AddCommand(
		install.NewCmd(c.out),
		script.NewCmd(c.out),
	)

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}

}
