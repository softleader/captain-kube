package app

import (
	"fmt"
	"github.com/spf13/cobra"
)

const ctxHelp = `Switch between Captain-Kubes back and forth

	ctx                    : list the contexts
	ctx                    : list the contexts with args
	ctx <NAME>             : switch to context <NAME>
	ctx -                  : switch to the previous context
	ctx -d <NAME>          : delete context <NAME> ('.' for current-context)
	ctx -a <NAME> <ARGS..> : add context <NAME> with <ARGS...>

	參數的讀取順序為: 當前 flags > ctx > os.Lookup
`

type ctxCmd struct {
	add    string
	delete string
	args   []string
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
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
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

	// TODO:
	// os.Lookup("SL_PLUGIN_MOUNT") 可以取得從 slctl 傳入的 plugin 資料儲存位置
	// 如果沒有發現 $SL_PLUGIN_MOUNT 這個 command 應該就要中斷
	// 反之 ctx 就可以儲存在那 $SL_PLUGIN_MOUNT 目錄下

	return nil
}
