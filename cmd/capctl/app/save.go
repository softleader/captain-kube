package app

import (
	"errors"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/dockerd"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/spf13/cobra"
	"path/filepath"
)

const (
	saveHelp = `匯出 Helm Chart 中的 image
`
)

type saveCmd struct {
	output string
	diff   bool
	charts []string
}

func newSaveCmd() *cobra.Command {
	c := saveCmd{}

	cmd := &cobra.Command{
		Use:   "save CHART",
		Short: "save images of a helm-chart",
		Long:  usage(saveHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			if c.charts = args; len(c.charts) == 0 {
				return errors.New("chart path is required")
			}
			return c.run()
		},
	}

	f := cmd.Flags()
	f.StringVarP(&c.output, "output", "o", c.output, "location of saved file")

	cmd.MarkFlagRequired("output")

	return cmd
}

func (c *saveCmd) run() error {
	expanded, err := homedir.Expand(c.output)
	if err != nil {
		expanded = c.output
	}
	abs, err := filepath.Abs(expanded)
	if err != nil {
		return err
	}
	tpls, err := chart.LoadArchive(logrus.StandardLogger(), abs)
	if err != nil {
		return err
	}

	var allImages []*chart.Image

	for _, images := range tpls {
		allImages = append(allImages, images...)
	}

	return dockerd.Save(logrus.StandardLogger(), allImages, c.output)
}
