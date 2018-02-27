package main

import (
	"github.com/softleader/captain-kube/cmd"
	"github.com/spf13/cobra"
)


func main() {

	var name string
	var version string
	var path, oldTagDomain, newTagDomain string

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

	/** 解壓縮 */
	var cmdExtract = cmd.Extract()
	/** 重新 tag image 後 push */
	var cmdRetag = cmd.Retag(path, oldTagDomain, newTagDomain)


	var rootCmd = &cobra.Command{Use: "ck"}
	rootCmd.AddCommand(cmdInstall, cmdUninstall, cmdExtract, cmdRetag)
	cmdInstall.AddCommand(cmdInstallByHelm, cmdInstallByICP)
	cmdUninstall.AddCommand(cmdUninstallByHelm, cmdUninstallByICP)

	rootCmd.Execute()

}
