package app

import (
	"errors"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/ctx"
	"github.com/softleader/captain-kube/pkg/dockerd"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	pb "github.com/softleader/captain-kube/pkg/proto"
	"github.com/spf13/cobra"
	"path/filepath"
)

const (
	pullHelp = `拉取 helm-chart 中的 image

	$ {{.}} pull CHART...

如果 registry 需要登入, 可以傳入 '--reg-*' 開頭的 flags 指定 docker registry 的認證資訊

	$ {{.}} pull CHART... --reg-user ME --reg-pass SECRET

如果需要在 pull 前修改 values.yaml 中任何參數, 可以傳入 '--set key1=val1,key2=val2'

	$ {{.}} pull CHART... --set ingress.enabled=true
`
)

type pullCmd struct {
	charts       []string
	registryAuth *ctx.RegistryAuth // docker registry auth
	helmChart    *ctx.HelmChart
}

func newPullCmd() *cobra.Command {
	c := pullCmd{
		registryAuth: activeCtx.RegistryAuth,
		helmChart:    activeCtx.HelmChart,
	}

	cmd := &cobra.Command{
		Use:   "pull CHART...",
		Short: "拉取 helm-chart 中的 image",
		Long:  usage(pullHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			if c.charts = args; len(c.charts) == 0 {
				return errors.New("chart path is required")
			}
			return c.run()
		},
	}

	f := cmd.Flags()
	c.registryAuth.AddFlags(f)
	c.helmChart.AddFlags(f)

	return cmd
}

func (c *pullCmd) run() error {
	for _, chart := range c.charts {
		if err := c.pull(chart); err != nil {
			return err
		}
	}
	return nil
}

func (c *pullCmd) pull(path string) error {
	expanded, err := homedir.Expand(path)
	if err != nil {
		expanded = path
	}
	abs, err := filepath.Abs(expanded)
	if err != nil {
		return err
	}
	tpls, err := chart.LoadArchive(logrus.StandardLogger(), abs, c.helmChart.Set...)
	if err != nil {
		return err
	}
	return dockerd.PullFromTemplates(logrus.StandardLogger(), tpls, &pb.RegistryAuth{
		Username: c.registryAuth.Username,
		Password: c.registryAuth.Password,
	})
}
