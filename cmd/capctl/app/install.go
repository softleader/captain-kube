package app

import (
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/ctx"
	"github.com/softleader/captain-kube/pkg/dockerd"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/spf13/cobra"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type installCmd struct {
	pull      bool
	sync      bool
	namespace string
	charts    []string

	retag        *ctx.ReTag
	registryAuth *ctx.RegistryAuth // docker registry auth
	helmTiller   *ctx.HelmTiller   // helm tiller
	endpoint     *ctx.Endpoint     // captain 的 endpoint ip
}

func newInstallCmd(activeCtx *ctx.Context) *cobra.Command {
	c := installCmd{
		namespace:    "default",
		retag:        activeCtx.ReTag,
		endpoint:     activeCtx.Endpoint,
		registryAuth: activeCtx.RegistryAuth,
		helmTiller:   activeCtx.HelmTiller,
	}

	cmd := &cobra.Command{
		Use:   "install [CHART...]",
		Short: "install helm-chart",
		Long:  "install helm chart",
		RunE: func(cmd *cobra.Command, args []string) error {
			if c.charts = args; len(c.charts) == 0 {
				return errors.New("chart path is required")
			}
			// do some validation check
			if err := c.endpoint.Validate(); err != nil {
				return err
			}
			// apply some default value
			if te := strings.TrimSpace(c.helmTiller.Endpoint); len(te) == 0 {
				c.helmTiller.Endpoint = c.endpoint.Host
			}
			return c.run()
		},
	}

	f := cmd.Flags()

	f.BoolVarP(&c.pull, "pull", "p", c.pull, "pull images in Chart")
	f.BoolVarP(&c.sync, "sync", "s", c.sync, "re-tag images & sync to all kubernetes nodes")

	f.StringVarP(&c.namespace, "namespace", "n", c.namespace, "specify the namespace of gcp, not available now")

	c.retag.AddFlags(f)
	c.endpoint.AddFlags(f)
	c.registryAuth.AddFlags(f)
	c.helmTiller.AddFlags(f)

	return cmd
}

func (c *installCmd) run() error {
	for _, chart := range c.charts {
		logrus.Println("### chart:", chart, "###")
		if err := runInstall(c, chart); err != nil {
			return err
		}
	}
	return nil
}

func runInstall(c *installCmd, path string) error {
	if expanded, err := homedir.Expand(path); err != nil {
		path = expanded
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadFile(abs)
	if err != nil {
		return err
	}

	request := proto.InstallChartRequest{
		Color:   settings.color,
		Timeout: settings.timeout,
		Verbose: settings.verbose,
		Chart: &proto.Chart{
			FileName: filepath.Base(abs),
			Content:  bytes,
			FileSize: int64(len(bytes)),
		},
		Pull: c.pull,
		Sync: c.sync,
		Retag: &proto.ReTag{
			From: c.retag.From,
			To:   c.retag.To,
		},
		Tiller: &proto.Tiller{
			Endpoint:          c.helmTiller.Endpoint,
			Username:          c.helmTiller.Username,
			Password:          c.helmTiller.Password,
			Account:           c.helmTiller.Account,
			SkipSslValidation: c.helmTiller.SkipSslValidation,
		},
		RegistryAuth: &proto.RegistryAuth{
			Username: c.registryAuth.Username,
			Password: c.registryAuth.Password,
		},
	}

	if err := dockerd.PullAndSync(logrus.StandardLogger(), &request); err != nil {
		return err
	}

	if err := captain.InstallChart(logrus.StandardLogger(), c.endpoint.String(), &request, settings.timeout); err != nil {
		return err
	}

	return nil
}
