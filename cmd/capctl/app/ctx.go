package app

import (
	"github.com/spf13/cobra"
)

const ctxHelp = `Switch between Captain-Kubes back and forth

	ctx ls                  : list the contexts
	ctx ls -a               : list the contexts with args
	ctx <NAME>              : switch to context <NAME>
	ctx -                   : switch to the previous context
	ctx rm <NAME>           : delete context <NAME> ('.' for current-context)
	ctx add <NAME> <ARGS..> : add context <NAME> with <ARGS...>
`

type ctxCmd struct {
}

func newCtxCmd() *cobra.Command {
	c := ctxCmd{}

	cmd := &cobra.Command{
		Use:   "ctx",
		Short: "switch between Captain-Kubes back and forth",
		Long:  ctxHelp,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.run()
		},
	}

	// f := cmd.Flags()

	return cmd
}

func (c *ctxCmd) run() error {
	return nil
}
