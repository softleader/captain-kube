package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"os/exec"
)

func Uninstall() *cobra.Command {

	var cmdUninstall = &cobra.Command{
		Use:   "uninstall [Helm/ICP/...]",
		Short: "uninstall Helm/ICP/... chart on k8s",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	return cmdUninstall
}


/** 用 helm 指令移除 chart */
func UninstallByHelm() *cobra.Command {

	var cmdUninstallByHelm = &cobra.Command{
		Use:   "helm [helm name]",
		Short: "Uninstall Helm chart on k8s",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			execCmd := exec.Command("sh", "-c", "helm delete --purge " + args[0])
			stdoutStderr := Output(execCmd.CombinedOutput())
			fmt.Printf("Finish uninstall  %s\n", stdoutStderr)

		},
	}
	return cmdUninstallByHelm
}


/** 用 ICP 的 bx 指令 uninstall chart */
func UninstallByICP(version string) *cobra.Command {

	var cmdUninstallByICP = &cobra.Command{
		Use:   "icp [chart name]",
		Short: "Uninstall Helm chart from ICP",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			unloadVersion := ""
			if version != "" {
				unloadVersion = " --version " + version
			}
			execCmd := exec.Command("sh", "-c", "bx pr delete-helm-chart --name " + args[0] + unloadVersion)

			stdoutStderr := Output(execCmd.CombinedOutput())
			fmt.Printf("Finish uninstall  %s\n", stdoutStderr)
		},
	}


	cmdUninstallByICP.Flags().StringVarP(&version, "version", "v", "", "which version of chart you want to unload.")

	return cmdUninstallByICP
}


