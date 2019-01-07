package app

import (
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/ctx"
	"github.com/softleader/captain-kube/pkg/dockerd"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/spf13/cobra"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type installCmd struct {
	pull bool // capctl 或 capui 是否要 pull image
	sync bool // 同步 image 到所有 node 上: 有 re-tag 時僅同步符合 re-tag 條件的 image; 無 re-tag 則同步全部
	//namespace    string
	charts       []string
	retag        *ctx.ReTag
	registryAuth *ctx.RegistryAuth // docker registry auth
	helmTiller   *ctx.HelmTiller   // helm tiller
	endpoint     *ctx.Endpoint     // captain 的 endpoint ip
}

func newInstallCmd() *cobra.Command {
	c := installCmd{
		//namespace:    "default",
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
	f.BoolVarP(&c.sync, "sync", "s", c.sync, "同步 image 到所有 node 上, 有 re-tag 則會同步 re-tag 之後的 image host")

	// f.StringVarP(&c.namespace, "namespace", "n", c.namespace, "specify the namespace of gcp, not available now")

	c.retag.AddFlags(f)
	c.endpoint.AddFlags(f)
	c.registryAuth.AddFlags(f)
	c.helmTiller.AddFlags(f)

	return cmd
}

func (c *installCmd) run() error {
	for _, chart := range c.charts {
		logrus.Printf("Installing helm chart: %s", chart)
		if err := runInstall(c, chart); err != nil {
			return err
		}
		logrus.Printf("Successfully installed chart to %q", c.helmTiller.Endpoint)
	}
	return nil
}

func runInstall(c *installCmd, path string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	expanded, err := homedir.Expand(path)
	if err != nil {
		path = expanded
	}
	bytes, err := ioutil.ReadFile(expanded)
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

	var tpls chart.Templates

	if c.pull {
		if tpls == nil {
			if tpls, err = chart.LoadBytes(logrus.StandardLogger(), request.Chart.Content); err != nil {
				return err
			}
		}
		if err := dockerd.PullFromTemplates(logrus.StandardLogger(), tpls, request.RegistryAuth); err != nil {
			return err
		}
	}

	if len(c.retag.From) > 0 && len(c.retag.To) > 0 {
		if tpls == nil {
			if tpls, err = chart.LoadBytes(logrus.StandardLogger(), request.Chart.Content); err != nil {
				return err
			}
		}
		if err := dockerd.ReTagFromTemplates(logrus.StandardLogger(), tpls, request.Retag, request.RegistryAuth); err != nil {
			return err
		}
	}

	if err := captain.InstallChart(logrus.StandardLogger(), c.endpoint.String(), &request, settings.timeout); err != nil {
		return err
	}

	return nil
}
