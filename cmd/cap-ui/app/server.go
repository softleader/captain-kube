package app

import (
	"fmt"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/capui"
	"github.com/softleader/captain-kube/pkg/ctx"
	"github.com/softleader/captain-kube/pkg/env"
	"github.com/spf13/cobra"
)

type capUiCmd struct {
	uiPort       int
	DefaultValue DefaultValue      // 畫面選ss項預設值
	RegistryAuth *ctx.RegistryAuth // docker registry 授權
	Tiller       *ctx.HelmTiller   // helm tiller 參數
	Endpoint     *ctx.Endpoint     // captain 的 endpoint ip
	EndpointStr  string            // endpoint.String() 由系統處理
}

type DefaultValue struct {
	Plaform   string // 平台(Google/ICP)
	Namespace string
	ReTag     *ctx.ReTag
}

func NewCapUiCommand() (cmd *cobra.Command) {
	var verbose bool
	envCtx := ctx.NewContextFromEnv()
	c := capUiCmd{
		DefaultValue: DefaultValue{
			Plaform:   env.Lookup(capui.EnvPlaform, capui.DefaultPlaform),
			Namespace: env.Lookup(capui.EnvNamespace, capui.DefaultNamespace),
			ReTag:     envCtx.ReTag,
		},

		RegistryAuth: envCtx.RegistryAuth,
		Tiller:       envCtx.HelmTiller,
		Endpoint:     envCtx.Endpoint,
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

	c.DefaultValue.ReTag.AddFlags(f)
	c.RegistryAuth.AddFlags(f)
	c.Tiller.AddFlags(f)
	c.Endpoint.AddFlags(f)

	return
}

func (cmd *capUiCmd) run() error {
	server := NewCapUiServer(cmd)
	return server.Run(fmt.Sprintf(":%v", cmd.uiPort))
}
