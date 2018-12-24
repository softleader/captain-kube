package cmd

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
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
	log            *logrus.Logger
	pull           bool
	sync           bool
	namespace      string
	sourceRegistry string
	registry       string
	charts         []string
	verbose        bool
	timeout        int64

	registryAuthUsername string // docker registry 的帳號
	registryAuthPassword string // docker registry 的密碼

	tillerEndpoint          string // helm tiller 的 ip, 若沒改預設為 endpoint
	tillerUsername          string // helm tiller 的使用者
	tillerPassword          string // helm tiller 的密碼
	tillerAccount           string // helm tiller 的帳號
	tillerSkipSslValidation bool

	endpoint     string // captain 的 endpoint ip
	endpointPort int    // captain 的 endpoint port
}

func NewInstallCmd(log *logrus.Logger, verbose bool) *cobra.Command {
	c := installCmd{
		log:       log,
		verbose:   verbose,
		namespace: "default",
		timeout:   300,

		registryAuthUsername: env.Lookup(captain.EnvRegistryAuthUsername, captain.DefaultRegistryAuthUsername),
		registryAuthPassword: env.Lookup(captain.EnvRegistryAuthPassword, captain.DefaultRegistryAuthPassword),

		tillerEndpoint:          env.Lookup(captain.EnvTillerEndpoint, ""),
		tillerUsername:          env.Lookup(captain.EnvTillerUsername, captain.DefaultTillerUsername),
		tillerPassword:          env.Lookup(captain.EnvTillerPassword, captain.DefaultTillerPassword),
		tillerAccount:           env.Lookup(captain.EnvTillerAccount, captain.DefaultTillerAccount),
		tillerSkipSslValidation: env.LookupBool(captain.EnvTillerSkipSslValidation, captain.DefaultTillerSkipSslValidation),

		endpointPort: captain.DefaultPort,
	}

	cmd := &cobra.Command{
		Use:   "install [CHART...]",
		Short: "install /path/to/chart.tgz",
		Long:  "install helm chart",
		RunE: func(cmd *cobra.Command, args []string) error {
			if c.charts = args; len(c.charts) == 0 {
				return errors.New("chart path is required")
			}
			// do some validation check
			if e := strings.TrimSpace(c.endpoint); len(e) == 0 {
				return errors.New("endpoint is required")
			}
			// apply some default value
			if te := strings.TrimSpace(c.tillerEndpoint); len(te) == 0 {
				c.tillerEndpoint = c.endpoint
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

	f.Int64Var(&c.timeout, "timeout", c.timeout, "seconds of captain run timeout")

	f.StringVar(&c.registryAuthUsername, "reg-user", c.registryAuthUsername, "specify the registryAuthUsername, reqiured when Pull&Sync")
	f.StringVar(&c.registryAuthPassword, "reg-pass", c.registryAuthPassword, "specify the registryAuthPassword, reqiured when Pull&Sync")

	f.StringVar(&c.tillerEndpoint, "tiller", c.tillerEndpoint, "specify the endpoint of helm tiller")
	f.StringVar(&c.tillerUsername, "tiller-user", c.tillerUsername, "specify the username of helm tiller")
	f.StringVar(&c.tillerPassword, "tiller-pass", c.tillerPassword, "specify the password of helm tiller")
	f.StringVar(&c.tillerAccount, "tiller-account", c.tillerAccount, "specify the account of helm tiller")
	f.BoolVar(&c.tillerSkipSslValidation, "tiller-skip-ssl", c.tillerSkipSslValidation, "specify skip ssl validation of helm tiller")

	f.StringVarP(&c.endpoint, "endpoint", "e", "", "specify the captain endpoint")
	f.IntVar(&c.endpointPort, "endpoint-port", captain.DefaultPort, "specify the port of captain endpoint")

	return cmd
}

func (c *installCmd) run() error {
	for _, chart := range c.charts {
		if err := run(c, chart); err != nil {
			return err
		}
	}
	return nil
}

func run(c *installCmd, path string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadFile(abs)
	if err != nil {
		return err
	}

	request := proto.InstallChartRequest{
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
		Timeout: c.timeout,
	}

	if err := dockerd.PullAndSync(c.log, &request); err != nil {
		return err
	}

	if err := captain.InstallChart(c.log, fmt.Sprintf("%s:%v", c.endpoint, c.endpointPort), &request, c.timeout); err != nil {
		return err
	}

	return nil
}
