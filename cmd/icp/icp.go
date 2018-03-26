package icp

import (
	"github.com/spf13/cobra"
	"os/exec"
	"fmt"
	"github.com/softleader/captain-kube/log"
	"strings"
)

/** 用 ICP 的 bx 指令 install chart */
func Install() (cmd *cobra.Command) {
	var verbose bool
	var options string
	cmd = &cobra.Command{
		Use:   "icp <Chart archive (.tgz)>",
		Short: "Install charts archive to IBM Cloud Private",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			commands := []string{"bx pr load-helm-chart --archive", args[0]}
			if options != "" {
				commands = append(commands, options)
			}
			c := strings.Join(commands, " ")
			if verbose {
				log.Command(c)
			}
			execCmd := exec.Command("sh", "-c", c)
			log.Output(execCmd.CombinedOutput())
			fmt.Printf("Installing %s, done.\n", args[0])
		},
	}
	cmd.Flags().StringVarP(&options, "options", "o", "", "Passing more options to underlying command")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Make the operation more talkative")
	return
}

/** 用 ICP 的 bx 指令 uninstall chart */
func Uninstall() (cmd *cobra.Command) {
	var verbose bool
	var options, version string
	cmd = &cobra.Command{
		Use:   "icp <chart name>",
		Short: "Uninstall charts archive to IBM Private Cloud",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			commands := []string{"bx pr delete-helm-chart --name ", args[0]}
			if version != "" {
				commands = append(commands, "--version", version)
			}
			if options != "" {
				commands = append(commands, options)
			}
			c := strings.Join(commands, " ")
			if verbose {
				log.Command(c)
			}
			execCmd := exec.Command("sh", "-c", c)
			log.Output(execCmd.CombinedOutput())
			fmt.Printf("Uninstalling %s, done.\n", args[0])
		},
	}
	cmd.Flags().StringVarP(&version, "version", "", "", "The Version of chart to unload.")
	cmd.Flags().StringVarP(&options, "options", "o", "", "Passing more options to underlying command")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Make the operation more talkative")
	return
}
