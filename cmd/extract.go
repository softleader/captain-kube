package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"os/exec"
)


func Extract() *cobra.Command {

	var cmdInstall = &cobra.Command{
		Use:   "tar [.tar.gz/.tar.tgz/.tar.Z/.tgz name]",
		Short: "tar .tar.gz/.tar.tgz/.tar.Z/.tgz file",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			execCmd := exec.Command("sh", "-c", "tar zxvf " + args[0])
			stdoutStderr := Output(execCmd.CombinedOutput())
			fmt.Printf("Finish extract  %s\n", stdoutStderr)
		},
	}
	return cmdInstall
}
