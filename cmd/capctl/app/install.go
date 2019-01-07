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

const (
	installHelp = `上傳一或多個 Helm Chart 至 Captain-Kube

使用 '--endpoint' 指定上傳的 Captain Endpoint

	$ {{.}} install CHART... -e CAPTAIN_ENDPOINT

若 Helm Tiller Server 不在 Captain-Kube 環境中, 可以傳入 '--tiller*' 開頭的 flag 設定 Tiller 相關資訊

	$ {{.}} install CHART... -e CAPTAIN_ENDPOINT --tiller TILLER_IP
	$ {{.}} install CHART... -e CAPTAIN_ENDPOINT --tiller TILLER_IP --tiller-skip-ssl=false

在上傳 Chart 之前, 支援以下 Pre-Procedures:

'--pull' : 拉取 Chart 中的 image
'--retag-from' 及 '--retag-to' : 將 Chart 中的 image tag 成指定 host 並推入該 docker registry

	$ {{.}} install CHART... -e CAPTAIN_ENDPOINT -p
	$ {{.}} install CHART... -e CAPTAIN_ENDPOINT -p -f hub.softleader.com.tw -t client-registry:5000

在上傳 Chart 之後, 支援以下 Post-Procedures:

'--sync' : 自動同步 image 到所有 kubernetes worker nodes, 如果當次上傳也有 re-tag 則會同步 re-tag 之後的 image

	$ {{.}} install CHART... -e CAPTAIN_ENDPOINT -s
	$ {{.}} install CHART... -e CAPTAIN_ENDPOINT -s -f hub.softleader.com.tw -t client-registry:5000

Pre-Procedures 跟 Post-Procedures 均可混合使用
亦可結合 '--reg-*' 開頭的 flags 指定 docker registry 的認證資訊

	$ {{.}} install CHART... -e CAPTAIN_ENDPOINT -ps --reg-pass SECRET
`
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
		Use:   "install CHART...",
		Short: "install helm-chart",
		Long:  usage(installHelp),
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
		if err := c.install(chart); err != nil {
			return err
		}
		logrus.Printf("Successfully installed chart to %q", c.helmTiller.Endpoint)
	}
	return nil
}

func (c *installCmd) install(path string) error {
	expanded, err := homedir.Expand(path)
	if err != nil {
		path = expanded
	}
	abs, err := filepath.Abs(expanded)
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadFile(abs)
	if err != nil {
		return err
	}

	request := proto.InstallChartRequest{
		Color:   settings.Color,
		Timeout: settings.Timeout,
		Verbose: settings.Verbose,
		Chart: &proto.Chart{
			FileName: filepath.Base(expanded),
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
			if tpls, err = chart.LoadArchiveBytes(logrus.StandardLogger(), request.Chart.FileName, request.Chart.Content); err != nil {
				return err
			}
		}
		if err := dockerd.PullFromTemplates(logrus.StandardLogger(), tpls, request.RegistryAuth); err != nil {
			return err
		}
	}

	if len(c.retag.From) > 0 && len(c.retag.To) > 0 {
		if tpls == nil {
			if tpls, err = chart.LoadArchiveBytes(logrus.StandardLogger(), request.Chart.FileName, request.Chart.Content); err != nil {
				return err
			}
		}
		if err := dockerd.ReTagFromTemplates(logrus.StandardLogger(), tpls, request.Retag, request.RegistryAuth); err != nil {
			return err
		}
	}

	if err := captain.InstallChart(logrus.StandardLogger(), c.endpoint.String(), &request, settings.Timeout); err != nil {
		return err
	}

	return nil
}
