package icp

import (
	"github.com/spf13/cobra"
	"os/exec"
	"fmt"
	"github.com/softleader/captain-kube/logs"
)

/** 用 ICP 的 bx 指令 install chart */
func Install() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "icp <Chart archive (.tgz)>",
		Short: "Install charts archive to IBM Private Cloud",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			execCmd := exec.Command("sh", "-c", "bx pr load-helm-chart --archive "+args[0])
			stdoutStderr := logs.Output(execCmd.CombinedOutput())
			fmt.Printf("Finish install  %s\n", stdoutStderr)
		},
	}
	return
}

/** 用 ICP 的 bx 指令 uninstall chart */
func Uninstall(version string) (cmd *cobra.Command) {

	cmd = &cobra.Command{
		Use:   "icp <chart name>",
		Short: "Uninstall charts archive to IBM Private Cloud",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			unloadVersion := ""
			if version != "" {
				unloadVersion = " --version " + version
			}
			execCmd := exec.Command("sh", "-c", "bx pr delete-helm-chart --name "+args[0]+unloadVersion)

			stdoutStderr := logs.Output(execCmd.CombinedOutput())
			fmt.Printf("Finish uninstall  %s\n", stdoutStderr)
		},
	}

	cmd.Flags().StringVarP(&version, "version", "v", "", "which version of chart you want to unload.")
	return
}
