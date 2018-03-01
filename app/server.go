package app

import (
	"github.com/spf13/cobra"
	"github.com/softleader/captain-kube/cmd/tar"
	"github.com/softleader/captain-kube/cmd/helm"
	"github.com/softleader/captain-kube/cmd/icp"
	"github.com/softleader/captain-kube/cmd/image"
)

const componentCaptainKube = "ckube"

func NewCaptainKubeCommand() (cmd *cobra.Command) {

	var name string
	var version string
	var path, oldTagDomain, newTagDomain string

	cmd = &cobra.Command{Use: componentCaptainKube}
	cmd.AddCommand(
		/** 安裝 **/
		install(name),
		/** 反安裝 **/
		uninstall(version),
		/** 解壓縮 */
		tar.Extract(),
		/** 重新 tag image 後 push */
		image.Retag(path, oldTagDomain, newTagDomain))

	return
}

func install(name string) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "install [Helm/ICP/...]",
		Short: "Install Helm/ICP/... chart on k8s",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	cmd.AddCommand(
		/** 利用 helm 指令 uninstall */
		helm.Install(name),
		/** 利用 ICP 的 bx 指令 helm */
		icp.Install())
	return
}

func uninstall(version string) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "uninstall [Helm/ICP/...]",
		Short: "uninstall Helm/ICP/... chart on k8s",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	cmd.AddCommand(
		/** 利用 helm 指令 uninstall */
		helm.Uninstall(),
		/** 利用 ICP 的 bx 指令 uninstall */
		icp.Uninstall(version))
	return
}
