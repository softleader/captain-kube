package helm

import (
	"github.com/spf13/cobra"
	"os/exec"
	"fmt"
	"github.com/softleader/captain-kube/log"
	"strings"
)

/** 用 helm 指令 install chart */
func Install() (cmd *cobra.Command) {
	var verbose bool
	var options, nickname string
	cmd = &cobra.Command{
		Use:   "helm <Chart directory>",
		Short: "Install charts to pure Kubernetes",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			commands := []string{"helm install", args[0], "-n", nickname}
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
	cmd.Flags().StringVarP(&nickname, "name", "n", "", "Helm charm Nickname (required)")
	cmd.Flags().StringVarP(&options, "options", "o", "", "Passing more options to underlying command")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Make the operation more talkative")
	cmd.MarkFlagRequired("name")
	return
}

/** 用 helm 指令移除 chart */
func Uninstall() (cmd *cobra.Command) {
	var verbose bool
	var options string
	cmd = &cobra.Command{
		Use:   "helm <nickname>",
		Short: "Uninstall charts from pure Kubernetes",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			commands := []string{"helm delete --purge", args[0]}
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
	cmd.Flags().StringVarP(&options, "options", "o", "", "Passing more options to underlying command")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Make the operation more talkative")
	return
}
