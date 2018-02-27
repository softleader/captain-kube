package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"os/exec"
)


func Install() *cobra.Command {

	var cmdInstall = &cobra.Command{
		Use:   "install [Helm/ICP/...]",
		Short: "Install Helm/ICP/... chart on k8s",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	return cmdInstall
}

/** 用 helm 指令 install chart */
func InstallByHelm(name string) *cobra.Command {

	var cmdInstallByHelm = &cobra.Command{
		Use:   "helm [chart file]",
		Short: "Install Helm chart file on k8s",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			execCmd := exec.Command("sh", "-c", "helm install " + args[0] + " -n " + name)
			stdoutStderr := Output(execCmd.CombinedOutput())
			fmt.Printf("Finish install  %s\n", stdoutStderr)

		},
	}
	cmdInstallByHelm.Flags().StringVarP(&name, "name", "n", "", "helm charm nickname (required)")
	cmdInstallByHelm.MarkFlagRequired("name")

	return cmdInstallByHelm
}

/** 用 ICP 的 bx 指令 install chart */
func InstallByICP() *cobra.Command {

	var cmdInstallByICP = &cobra.Command{
		Use:   "icp [chart archive (.tgz)]",
		Short: "Install Helm chart archive on ICP",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			execCmd := exec.Command("sh", "-c", "bx pr load-helm-chart --archive " + args[0])
			stdoutStderr := Output(execCmd.CombinedOutput())
			fmt.Printf("Finish install  %s\n", stdoutStderr)
		},
	}

	return cmdInstallByICP
}
