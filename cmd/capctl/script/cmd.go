package script

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
)

type scriptCmd struct {
	out       io.Writer
	tags      string
	chartPath string
}

func NewCmd(out io.Writer) *cobra.Command {
	c := scriptCmd{
		out: out,
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
