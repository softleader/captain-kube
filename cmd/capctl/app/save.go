package app

import (
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/dockerd"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/spf13/cobra"
	"path/filepath"
)

const (
	saveHelp = `匯出一個或多個 Helm Chart 中的 image

傳入 '--output' 指定儲存的檔案路徑, docker 預設的儲存檔案格式為 tarball (.tar)

	$ {{.}} save CHART... -o OUTPUT.tar

傳入 '--force' 可以強制複寫已存在的 output 檔案

	$ {{.}} save CHART... -o OUTPUT.tar -f
`
)

type saveCmd struct {
	output string
	force  bool
	charts []string
}

func newSaveCmd() *cobra.Command {
	c := saveCmd{}

	cmd := &cobra.Command{
		Use:   "save CHART...",
		Short: "save images of helm-charts",
		Long:  usage(saveHelp),
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c.charts = args
			return c.run()
		},
	}

	f := cmd.Flags()
	f.StringVarP(&c.output, "output", "o", c.output, "location of saved file")
	f.BoolVarP(&c.force, "force", "f", false, "force to delete output file if exist")

	cmd.MarkFlagRequired("output")

	return cmd
}

func (c *saveCmd) run() error {
	var allImages []*chart.Image
	for _, path := range c.charts {
		expanded, err := homedir.Expand(path)
		if err != nil {
			expanded = c.output
		}
		abs, err := filepath.Abs(expanded)
		if err != nil {
			return err
		}
		logrus.Printf("Collecting images from: %s\n", abs)
		tpls, err := chart.LoadArchive(logrus.StandardLogger(), abs)
		if err != nil {
			return err
		}

		for tpl, images := range tpls {
			logrus.Debugf("detecting source: %s\n", tpl)
			for _, image := range images {
				logrus.Println(image)
				allImages = append(allImages, image)
			}
		}
	}
	return dockerd.Save(logrus.StandardLogger(), allImages, c.output, c.force)
}
