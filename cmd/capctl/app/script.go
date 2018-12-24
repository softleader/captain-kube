package app

import (
	"fmt"
	"github.com/spf13/cobra"
)

type scriptCmd struct {
	tags      string
	chartPath string
}

func newScriptCmd() *cobra.Command {
	c := scriptCmd{
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
