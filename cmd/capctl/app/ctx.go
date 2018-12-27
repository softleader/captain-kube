package app

import (
	"encoding/json"
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/ctx"
	"github.com/spf13/cobra"
)

const (
	ctxHelp = `Switch between captain-kubes back and forth

	ctx                       : 列出所有 context
	ctx --width 0             : 列出所有 context 並顯示完整的 args
	ctx <NAME>                : 切換 context 到 <NAME>
	ctx -                     : 切換到前一個 context
	ctx -d <NAME>             : 刪除 context <NAME> ('.' 為當前的 context)
	ctx -a <NAME> -- <ARGS..> : 新增 context <NAME>

參數的讀取順序為: 當前 flags > ctx > os.LookupEnv
`
)

type ctxCmd struct {
	width  uint
	add    string
	delete string
	args   []string
	ctxs   *ctx.Contexts
}

func newCtxCmd(ctxs *ctx.Contexts) *cobra.Command {
	c := ctxCmd{}

	cmd := &cobra.Command{
		Use:   "ctx",
		Short: "switch between captain-kubes back and forth",
		Long:  ctxHelp,
		Args: func(cmd *cobra.Command, args []string) error {
			if ctxs == ctx.PlainContexts {
				return ctx.ErrMountVolumeNotExist
			}
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
			c.ctxs = ctxs
			c.args = args
			return c.run()
		},
	}

	f := cmd.Flags()
	f.StringVarP(&c.add, "add", "a", "", "add context <NAME> with <ARGS...>")
	f.StringVarP(&c.delete, "delete", "d", "", "delete context <NAME> ('.' for current-context)")
	f.UintVar(&c.width, "width", 100, "maximum allowed width for listing context args")

	return cmd
}

func (c *ctxCmd) run() error {
	if len(c.add) > 0 {
		return c.ctxs.Add(c.add, c.args)
	}
	if len(c.delete) > 0 {
		return c.ctxs.Delete(c.delete)
	}
	if len(c.args) > 0 {
		return c.ctxs.Switch(c.args[0])
	}

	table := uitable.New()
	table.AddRow("", "CONTEXT", "ARGS")
	table.MaxColWidth = c.width
	for name, ctx := range c.ctxs.Contexts {
		prefix := " "
		if name == c.ctxs.Active {
			prefix = ">"
		} else if name == c.ctxs.Previous {
			prefix = "-"
		}
		args, err := json.Marshal(ctx)
		if err != nil {
			return err
		}
		table.AddRow(prefix, name, fmt.Sprintf("%+v", string(args)))
	}
	logrus.Println(table)
	return nil
}
