package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"os/exec"
	"log"
)


type Args struct {
	Source   string
	Dest     string
}

func Install() *cobra.Command {

	var cmdInstall = &cobra.Command{
		Use:   "install [Helm/ICP/...]",
		Short: "Install Helm/ICP/... chart on k8s",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			fmt.Printf("Finish install \n")
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
			stdoutStderr := output(execCmd.CombinedOutput())
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
			stdoutStderr := output(execCmd.CombinedOutput())
			fmt.Printf("Finish install  %s\n", stdoutStderr)
		},
	}

	return cmdInstallByICP
}


func Uninstall() *cobra.Command {

	var cmdUninstall = &cobra.Command{
		Use:   "uninstall [Helm/ICP/...]",
		Short: "uninstall Helm/ICP/... chart on k8s",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			fmt.Printf("Finish uninstall \n")
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
			stdoutStderr := output(execCmd.CombinedOutput())
			fmt.Printf("Finish install  %s\n", stdoutStderr)

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

			stdoutStderr := output(execCmd.CombinedOutput())
			fmt.Printf("Finish install  %s\n", stdoutStderr)
		},
	}


	cmdUninstallByICP.Flags().StringVarP(&version, "version", "v", "", "which version of chart you want to unload.")

	return cmdUninstallByICP
}


func output(stdoutStderr []byte, err error) []byte {
	if err != nil {
		fmt.Printf("%s\n", stdoutStderr)
		log.Fatal(err)
	}
	return stdoutStderr
}
