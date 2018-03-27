package tar

import (
	"github.com/spf13/cobra"
	"os/exec"
	"fmt"
	"github.com/softleader/captain-kube/log"
)

func Extract() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "tar [tar file]",
		Short: "Extract tar",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			execCmd := exec.Command("sh", "-c", "tar zxvf "+args[0])
			fmt.Println("Finish extract", args[0])
			log.Output(execCmd.CombinedOutput())
		},
	}
	return
}
