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
	cmd = &cobra.Command{Use: componentCaptainKube}
	cmd.AddCommand(
		/** 安裝 **/
		install(),
		/** 反安裝 **/
		uninstall(),
		/** 解壓縮 */
		tar.Extract(),
		/** 重新 tag image 後 push */
		image.Retag())
	return
}

func install() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "install <kind>",
		Short: "Install Helm Charts to Kubernetes",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	cmd.AddCommand(
		/** 利用 helm 指令 uninstall */
		helm.Install(),
		/** 利用 ICP 的 bx 指令 helm */
		icp.Install())
	return
}

func uninstall() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "uninstall <kind>",
		Short: "Uninstall Helm Charts from Kubernetes",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	cmd.AddCommand(
		/** 利用 helm 指令 uninstall */
		helm.Uninstall(),
		/** 利用 ICP 的 bx 指令 uninstall */
		icp.Uninstall())
	return
}
