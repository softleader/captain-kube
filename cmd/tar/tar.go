package tar

import (
	"github.com/spf13/cobra"
	"os/exec"
	"fmt"
	"github.com/softleader/captain-kube/log"
)

func Extract() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "tar [.tar.gz/.tar.tgz/.tar.Z/.tgz name]",
		Short: "tar .tar.gz/.tar.tgz/.tar.Z/.tgz file",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			execCmd := exec.Command("sh", "-c", "tar zxvf "+args[0])
			stdoutStderr := log.Output(execCmd.CombinedOutput())
			fmt.Printf("Finish extract  %s\n", stdoutStderr)
		},
	}
	return
}
