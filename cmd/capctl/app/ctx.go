package app

import (
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/ctx"
	"github.com/spf13/cobra"
)

const (
	ctxHelp = `Switch between Captain-Kubes back and forth

	ctx                    : list the contexts
	ctx                    : list the contexts with args
	ctx <NAME>             : switch to context <NAME>
	ctx -                  : switch to the previous context
	ctx -d <NAME>          : delete context <NAME> ('.' for current-context)
	ctx -a <NAME> <ARGS..> : add context <NAME> with <ARGS...>

	參數的讀取順序為: 當前 flags > ctx > os.LookupEnv
`
)

type ctxCmd struct {
	add    string
	delete string
	args   []string
	ctxs   *ctx.Contexts
}

func newCtxCmd() *cobra.Command {
	c := ctxCmd{}

	cmd := &cobra.Command{
		Use:   "ctx",
		Short: "switch between Captain-Kubes back and forth",
		Long:  ctxHelp,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(c.add) > 0 && len(c.delete) > 0 {
				return fmt.Errorf("can not add and delete at the same time")
			}
			if len(c.add) > 0 && len(args) == 0 {
				return fmt.Errorf("requires at least 1 argument to add context")
			}
			if len(c.delete) > 0 && len(args) > 0 {
				return fmt.Errorf("delete context does not accpet arguments")
			}
			if len(c.add) == 0 && len(c.delete) == 0 && len(args) > 1 {
				return fmt.Errorf("list/switch context only accpet max 1 argument")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			loaded, err := ctx.LoadContextsFromEnv(logrus.StandardLogger())
			if err != nil {
				return err
			}
			c.ctxs = loaded
			c.args = args
			return c.run()
		},
	}

	f := cmd.Flags()
	f.StringVarP(&c.add, "add", "a", "", "add context <NAME> with <ARGS...>")
	f.StringVarP(&c.delete, "delete", "d", "", "delete context <NAME> ('.' for current-context)")

	return cmd
}

func (c *ctxCmd) run() error {
	if len(c.add) > 0 {
		return c.ctxs.Add(c.args[0], c.args[1:])
	}
	if len(c.delete) > 0 {
		return c.ctxs.Delete(c.args[0])
	}
	if len(c.args) > 0 {
		return c.ctxs.Switch(c.args[0])
	}

	table := uitable.New()
	for k, ctx := range c.ctxs.Contexts {
		prefix := " "
		if k == c.ctxs.Active {
			prefix = ">"
		} else if k == c.ctxs.Previous {
			prefix = "-"
		}
		table.AddRow(fmt.Sprintf("%s %s", prefix, k), fmt.Sprintf("%+v", ctx))
	}
	logrus.Println(table)
	return nil
}
