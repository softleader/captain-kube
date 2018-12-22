package script

import (
	"fmt"
	"github.com/softleader/captain-kube/pkg/logger"
	"github.com/spf13/cobra"
)

type scriptCmd struct {
	log       *logger.Logger
	tags      string
	chartPath string
}

func NewCmd(log *logger.Logger) *cobra.Command {
	c := scriptCmd{
		log: log,
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
