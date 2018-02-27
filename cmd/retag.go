package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os/exec"
	"strings"
	"github.com/softleader/captain-kube/app"
)

func Retag(path, oldTagDomain, newTagDomain string) *cobra.Command {

	var cmdRetag = &cobra.Command{
		Use:   "retag [remote/local]",
		Short: "retag docker image",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			m := app.ParseYaml(path)
			p := &app.KeywordValues{}
			app.FindKeywordFromMap(m, "repository", p)

			/** remote retag 方式 */
			if args[0] == "remote" {
				for _, v := range p.MappingValues {
					if strings.Contains(v, oldTagDomain) {
						newTag := strings.Replace(v, oldTagDomain, newTagDomain, -1)
						execPullAndRetagCmd := exec.Command("sh", "-c", "docker pull " + v + " && docker tag " + v + " " + newTag)
						fmt.Printf("Finish retag \n %s", Output(execPullAndRetagCmd.CombinedOutput()))

						execPushCmd := exec.Command("sh", "-c", "docker push " + newTag)
						fmt.Printf("Finish push \n %s", Output(execPushCmd.CombinedOutput()))
					}
				}

			}

			//TODO
			/** local retag 方式 */
		},
	}

	cmdRetag.Flags().StringVarP(&path, "path", "p", "", "yaml/yml path which you want to parse. (required)")
	cmdRetag.Flags().StringVarP(&oldTagDomain, "oldTagDomain", "o", "", "which tag domain you want to reTag. (required)")
	cmdRetag.Flags().StringVarP(&newTagDomain, "newTagDomain", "n", "", "which new tag domain name you want. (required)")
	cmdRetag.MarkFlagRequired("path")
	cmdRetag.MarkFlagRequired("oldTagDomain")
	cmdRetag.MarkFlagRequired("newTagDomain")

	return cmdRetag
}