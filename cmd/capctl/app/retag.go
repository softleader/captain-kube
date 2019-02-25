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
	retagHelp = `重新 tag helm-chart 中的 image

傳入 '--retag-from' 及 '--retag-to' 指定 retag 的來源及目標, 只有符合 '--retag-from' 的 image 才會被 retag

	$ {{.}} retag CHART... -f hub.softleader.com.tw -t client-registry:5000

如果 registry 需要登入, 可以傳入 '--reg-*' 開頭的 flags 指定 docker registry 的認證資訊

	$ {{.}} retag CHART... -f hub.softleader.com.tw -t client-registry:5000 --reg-pass SECRET
`
)

type retagCmd struct {
	charts       []string
	registryAuth *ctx.RegistryAuth // docker registry auth
	retag        *ctx.ReTag
}

func newReTagCmd() *cobra.Command {
	c := retagCmd{
		registryAuth: activeCtx.RegistryAuth,
		retag:        activeCtx.ReTag,
	}

	cmd := &cobra.Command{
		Use:   "retag CHART...",
		Short: "重新 tag helm-chart 中的 image",
		Long:  usage(retagHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			if c.charts = args; len(c.charts) == 0 {
				return errors.New("chart path is required")
			}
			return c.run()
		},
	}

	f := cmd.Flags()
	c.registryAuth.AddFlags(f)
	c.retag.AddFlags(f)

	return cmd
}

func (c *retagCmd) run() error {
	for _, chart := range c.charts {
		if err := c.reTag(chart); err != nil {
			return err
		}
	}
	return nil
}

func (c *retagCmd) reTag(path string) error {
	expanded, err := homedir.Expand(path)
	if err != nil {
		expanded = path
	}
	abs, err := filepath.Abs(expanded)
	if err != nil {
		return err
	}
	tpls, err := chart.LoadArchive(logrus.StandardLogger(), abs)
	if err != nil {
		return err
	}
	return dockerd.ReTagFromTemplates(logrus.StandardLogger(), tpls, &pb.ReTag{
		From: c.retag.From,
		To:   c.retag.To,
	}, &pb.RegistryAuth{
		Username: c.registryAuth.Username,
		Password: c.registryAuth.Password,
	})
}
