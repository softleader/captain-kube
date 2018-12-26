package app

import (
	"fmt"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/capui"
	"github.com/softleader/captain-kube/pkg/env"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/spf13/cobra"
)

type capUiCmd struct {
	uiPort       int
	DefaultValue DefaultValue      // 畫面選ss項預設值
	RegistryAuth RegistryAuth      // docker registry 授權
	Tiller       Tiller            // helm tiller 參數
	Endpoint     *captain.Endpoint // captain 的 endpoint ip
	EndpointStr  string            // endpoint.String() 由系統處理
}

type DefaultValue struct {
	Plaform   string // 平台(Google/ICP)
	Namespace string
	ReTag     proto.ReTag
}

type RegistryAuth struct {
	Username string // docker registry 的帳號
	Password string // docker registry 的密碼
}

type Tiller struct {
	Endpoint          string // helm tiller 的 ip, 若沒改預設為 endpoint
	Username          string // helm tiller 的使用者
	Password          string // helm tiller 的密碼
	Account           string // helm tiller 的帳號
	SkipSslValidation bool
}

func NewCapUiCommand() (cmd *cobra.Command) {
	var verbose bool

	c := capUiCmd{
		DefaultValue: DefaultValue{
			Plaform:   env.Lookup(capui.EnvPlaform, capui.DefaultPlaform),
			Namespace: env.Lookup(capui.EnvNamespace, capui.DefaultNamespace),
			ReTag: proto.ReTag{
				From: env.Lookup(capui.EnvReTagFrom, capui.DefaultReTagFrom),
				To:   env.Lookup(capui.EnvReTagTo, capui.DefaultReTagTo),
			},
		},

		RegistryAuth: RegistryAuth{
			Username: env.Lookup(captain.EnvRegistryAuthUsername, captain.DefaultRegistryAuthUsername),
			Password: env.Lookup(captain.EnvRegistryAuthPassword, captain.DefaultRegistryAuthPassword),
		},

		Tiller: Tiller{
			Endpoint:          env.Lookup(captain.EnvTillerEndpoint, ""),
			Username:          env.Lookup(captain.EnvTillerUsername, captain.DefaultTillerUsername),
			Password:          env.Lookup(captain.EnvTillerPassword, captain.DefaultTillerPassword),
			Account:           env.Lookup(captain.EnvTillerAccount, captain.DefaultTillerAccount),
			SkipSslValidation: env.LookupBool(captain.EnvTillerSkipSslValidation, captain.DefaultTillerSkipSslValidation),
		},
	}

	cmd = &cobra.Command{
		Use:  "capui",
		Long: "capui is a web interface for captain",
		RunE: func(cmd *cobra.Command, args []string) error {
			logrus.SetOutput(colorable.NewColorableStdout()) // for windows color output
			if verbose {
				logrus.SetLevel(logrus.DebugLevel)
			}
			if err := c.Endpoint.Validate(); err != nil {
				return err
			} else {
				c.EndpointStr = c.Endpoint.String()
			}

			return c.run()
		},
	}

	f := cmd.Flags()
	f.BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
	f.IntVarP(&c.uiPort, "port", "p", 8080, "port of web ui serve port")

	f.StringVarP(&c.DefaultValue.Plaform, "plaform", "k", c.DefaultValue.Plaform, "default value of k8s plaform")
	f.StringVarP(&c.DefaultValue.Namespace, "namespace", "n", c.DefaultValue.Namespace, "default value of the namespace of gcp, not available now")
	f.StringVarP(&c.DefaultValue.ReTag.From, "retag-from", "f", c.DefaultValue.ReTag.From, "default value of the registryAuthUsername")
	f.StringVarP(&c.DefaultValue.ReTag.To, "retag-to", "t", c.DefaultValue.ReTag.To, "default value of the registryAuthPassword")

	f.StringVar(&c.RegistryAuth.Username, "reg-user", c.RegistryAuth.Username, "specify the registryAuthUsername")
	f.StringVar(&c.RegistryAuth.Password, "reg-pass", c.RegistryAuth.Password, "specify the registryAuthPassword")

	f.StringVar(&c.Tiller.Endpoint, "tiller", c.Tiller.Endpoint, "specify the endpoint of helm tiller")
	f.StringVar(&c.Tiller.Username, "tiller-user", c.Tiller.Username, "specify the username of helm tiller")
	f.StringVar(&c.Tiller.Password, "tiller-pass", c.Tiller.Password, "specify the password of helm tiller")
	f.StringVar(&c.Tiller.Account, "tiller-account", c.Tiller.Account, "specify the account of helm tiller")
	f.BoolVar(&c.Tiller.SkipSslValidation, "tiller-skip-ssl", c.Tiller.SkipSslValidation, "specify skip ssl validation of helm tiller")

	c.Endpoint = captain.AddEndpointFlags(f)

	return
}

func (cmd *capUiCmd) run() error {
	server := NewCapUiServer(cmd)
	return server.Run(fmt.Sprintf(":%v", cmd.uiPort))
}
