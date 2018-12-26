package app

import (
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/dockerd"
	"github.com/softleader/captain-kube/pkg/env"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/spf13/cobra"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type installCmd struct {
	pull           bool
	sync           bool
	namespace      string
	sourceRegistry string
	registry       string
	charts         []string

	registryAuthUsername string // docker registry 的帳號
	registryAuthPassword string // docker registry 的密碼

	tillerEndpoint          string // helm tiller 的 ip, 若沒改預設為 endpoint
	tillerUsername          string // helm tiller 的使用者
	tillerPassword          string // helm tiller 的密碼
	tillerAccount           string // helm tiller 的帳號
	tillerSkipSslValidation bool

	endpoint *captain.Endpoint // captain 的 endpoint ip
}

func newInstallCmd() *cobra.Command {
	c := installCmd{
		namespace: "default",

		registryAuthUsername: env.Lookup(captain.EnvRegistryAuthUsername, captain.DefaultRegistryAuthUsername),
		registryAuthPassword: env.Lookup(captain.EnvRegistryAuthPassword, captain.DefaultRegistryAuthPassword),

		tillerEndpoint:          env.Lookup(captain.EnvTillerEndpoint, ""),
		tillerUsername:          env.Lookup(captain.EnvTillerUsername, captain.DefaultTillerUsername),
		tillerPassword:          env.Lookup(captain.EnvTillerPassword, captain.DefaultTillerPassword),
		tillerAccount:           env.Lookup(captain.EnvTillerAccount, captain.DefaultTillerAccount),
		tillerSkipSslValidation: env.LookupBool(captain.EnvTillerSkipSslValidation, captain.DefaultTillerSkipSslValidation),
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
			if te := strings.TrimSpace(c.tillerEndpoint); len(te) == 0 {
				c.tillerEndpoint = c.endpoint.Host
			}
			return c.run()
		},
	}

	f := cmd.Flags()

	f.BoolVarP(&c.pull, "pull", "p", c.pull, "pull images in Chart")
	f.BoolVarP(&c.sync, "sync", "s", c.sync, "re-tag images & sync to all kubernetes nodes")

	f.StringVarP(&c.namespace, "namespace", "n", c.namespace, "specify the namespace of gcp, not available now")

	f.StringVarP(&c.sourceRegistry, "retag-from", "f", c.sourceRegistry, "specify the host of re-tag from, required when Sync")
	f.StringVarP(&c.registry, "retag-to", "t", c.registry, "specify the host of re-tag to, required when Sync")

	f.StringVar(&c.registryAuthUsername, "reg-user", c.registryAuthUsername, "specify the registryAuthUsername, reqiured when Pull&Sync")
	f.StringVar(&c.registryAuthPassword, "reg-pass", c.registryAuthPassword, "specify the registryAuthPassword, reqiured when Pull&Sync")

	f.StringVar(&c.tillerEndpoint, "tiller", c.tillerEndpoint, "specify the endpoint of helm tiller")
	f.StringVar(&c.tillerUsername, "tiller-user", c.tillerUsername, "specify the username of helm tiller")
	f.StringVar(&c.tillerPassword, "tiller-pass", c.tillerPassword, "specify the password of helm tiller")
	f.StringVar(&c.tillerAccount, "tiller-account", c.tillerAccount, "specify the account of helm tiller")
	f.BoolVar(&c.tillerSkipSslValidation, "tiller-skip-ssl", c.tillerSkipSslValidation, "specify skip ssl validation of helm tiller")

	c.endpoint = captain.AddEndpointFlags(f)

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
			From: c.sourceRegistry,
			To:   c.registry,
		},
		Tiller: &proto.Tiller{
			Endpoint:          c.tillerEndpoint,
			Username:          c.tillerUsername,
			Password:          c.tillerPassword,
			Account:           c.tillerAccount,
			SkipSslValidation: c.tillerSkipSslValidation,
		},
		RegistryAuth: &proto.RegistryAuth{
			Username: c.registryAuthUsername,
			Password: c.registryAuthPassword,
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
