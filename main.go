package main

import (
	"github.com/softleader/captain-kube/cmd"
	"github.com/spf13/cobra"
)


func main() {

	var name string
	var version string

	var cmdInstall = cmd.Install()

	/** 利用 helm install*/
	var cmdInstallByHelm = cmd.InstallByHelm(name)
	/** 利用 ICP 的 bx 指令 install*/
	var cmdInstallByICP = cmd.InstallByICP()


	var cmdUninstall = cmd.Uninstall()

	/** 利用 helm 指令 uninstall */
	var cmdUninstallByHelm = cmd.UninstallByHelm()
	/** 利用 ICP 的 bx 指令 uninstall */
	var cmdUninstallByICP = cmd.UninstallByICP(version)



	var rootCmd = &cobra.Command{Use: "ck"}
	rootCmd.AddCommand(cmdInstall, cmdUninstall  /*, cmdLoad, cmdUnload*/)
	cmdInstall.AddCommand(cmdInstallByHelm, cmdInstallByICP)
	cmdUninstall.AddCommand(cmdUninstallByHelm, cmdUninstallByICP)
	rootCmd.Execute()
}
