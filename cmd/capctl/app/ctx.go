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

	參數的讀取順序為: 當前 flags > ctx > os.Lookup
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

	// TODO:
	// os.Lookup("SL_PLUGIN_MOUNT") 可以取得從 slctl 傳入的 plugin 資料儲存位置
	// 如果沒有發現 $SL_PLUGIN_MOUNT 這個 command 應該就要中斷
	// 反之 ctx 就可以儲存在那 $SL_PLUGIN_MOUNT 目錄下

	return nil
}
