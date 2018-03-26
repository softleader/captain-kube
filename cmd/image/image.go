package image

import (
	"github.com/spf13/cobra"
	"strings"
	"os/exec"
	"fmt"
	"github.com/softleader/captain-kube/log"
	"github.com/softleader/captain-kube/charts"
)

func Retag() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "retag <origin>",
		Short: "ReTag docker image from origin",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	cmd.AddCommand(
		/** retag 自 remote */
		remote(),
		/** retag 自 local */
		local())
	return
}

func remote() (cmd *cobra.Command) {
	var m map[interface{}]interface{}
	var key charts.KeywordValues
	var verbose bool
	var options, path, oldTagDomain, newTagDomain string
	cmd = &cobra.Command{
		Use:   "remote",
		Short: "from remote",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			m = charts.Load(path)
			key = charts.KeywordValues{}
			charts.FindKeywordFromMap(m, "repository", &key)

			for _, v := range key.MappingValues {
				if strings.Contains(v, oldTagDomain) {
					newTag := strings.Replace(v, oldTagDomain, newTagDomain, -1)
					execPullAndRetagCmd := exec.Command("sh", "-c", "docker pull "+v+" && docker tag "+v+" "+newTag)
					fmt.Printf("Finish retag")
					log.Output(execPullAndRetagCmd.CombinedOutput())

					execPushCmd := exec.Command("sh", "-c", "docker push "+newTag)
					fmt.Printf("Finish push")
					log.Output(execPushCmd.CombinedOutput())
				}
			}
		},
	}
	cmd.Flags().StringVarP(&path, "path", "p", "", "yaml/yml path which you want to parse. (required)")
	cmd.Flags().StringVarP(&oldTagDomain, "oldTagDomain", "", "", "which tag domain you want to reTag. (required)")
	cmd.Flags().StringVarP(&newTagDomain, "newTagDomain", "", "", "which new tag domain name you want. (required)")
	cmd.MarkFlagRequired("path")
	cmd.MarkFlagRequired("oldTagDomain")
	cmd.MarkFlagRequired("newTagDomain")

	cmd.Flags().StringVarP(&options, "options", "o", "", "Passing more options to underlying command")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Make the operation more talkative")
	return
}

func local() (cmd *cobra.Command) {
	var verbose bool
	var options string
	cmd = &cobra.Command{
		Use:   "local",
		Short: "from local",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO
		},
	}
	cmd.Flags().StringVarP(&options, "options", "o", "", "Passing more options to underlying command")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Make the operation more talkative")
	return
}
