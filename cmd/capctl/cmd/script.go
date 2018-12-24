package cmd

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

type scriptCmd struct {
	log       *logrus.Logger
	verbose   bool
	tags      string
	chartPath string
}

func NewScriptCmd(log *logrus.Logger, verbose bool) *cobra.Command {
	c := scriptCmd{
		log:     log,
		verbose: verbose,
	}

	cmd := &cobra.Command{
		Use:   "script <prsl>",
		Short: "build script of helm chart",
		Long:  "build script of helm chart",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.run()
		},
	}

	//f := cmd.Flags()

	return cmd
}

func (c *scriptCmd) run() error {
	fmt.Println("run script", c)
	return nil
}
