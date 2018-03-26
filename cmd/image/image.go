package image

import (
	"github.com/spf13/cobra"
	"strings"
	"os/exec"
	"fmt"
	"github.com/softleader/captain-kube/logs"
	"github.com/softleader/captain-kube/charts"
)

func Retag(path, oldTagDomain, newTagDomain string) (cmd *cobra.Command) {

	cmd = &cobra.Command{
		Use:   "retag [remote/local]",
		Short: "ReTag docker image",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			m := charts.Load(path)
			p := &charts.KeywordValues{}
			charts.FindKeywordFromMap(m, "repository", p)

			/** remote retag 方式 */
			if args[0] == "remote" {
				for _, v := range p.MappingValues {
					if strings.Contains(v, oldTagDomain) {
						newTag := strings.Replace(v, oldTagDomain, newTagDomain, -1)
						execPullAndRetagCmd := exec.Command("sh", "-c", "docker pull "+v+" && docker tag "+v+" "+newTag)
						fmt.Printf("Finish retag \n %s", logs.Output(execPullAndRetagCmd.CombinedOutput()))

						execPushCmd := exec.Command("sh", "-c", "docker push "+newTag)
						fmt.Printf("Finish push \n %s", logs.Output(execPushCmd.CombinedOutput()))
					}
				}

			}

			//TODO
			/** local retag 方式 */
		},
	}

	cmd.Flags().StringVarP(&path, "path", "p", "", "yaml/yml path which you want to parse. (required)")
	cmd.Flags().StringVarP(&oldTagDomain, "oldTagDomain", "o", "", "which tag domain you want to reTag. (required)")
	cmd.Flags().StringVarP(&newTagDomain, "newTagDomain", "n", "", "which new tag domain name you want. (required)")
	cmd.MarkFlagRequired("path")
	cmd.MarkFlagRequired("oldTagDomain")
	cmd.MarkFlagRequired("newTagDomain")
	return
}
