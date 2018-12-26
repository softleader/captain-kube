package app

import (
	"fmt"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/cmd/cap-ui/app/server"
	"github.com/softleader/captain-kube/cmd/cap-ui/app/server/comm"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/spf13/cobra"
)

type capUiCmd struct {
	configPath string
	uiPort     int

	registryAuthUsername string // docker registry 的帳號
	registryAuthPassword string // docker registry 的密碼

	tillerEndpoint          string // helm tiller 的 ip, 若沒改預設為 endpoint
	tillerUsername          string // helm tiller 的使用者
	tillerPassword          string // helm tiller 的密碼
	tillerAccount           string // helm tiller 的帳號
	tillerSkipSslValidation bool

	endpoint *captain.Endpoint // captain 的 endpoint ip
}

func NewCapUiCommand() (cmd *cobra.Command) {
	var verbose bool
	c := capUiCmd{}
	cmd = &cobra.Command{
		Use:  "capui",
		Long: "capui is a web interface for captain",
		RunE: func(cmd *cobra.Command, args []string) error {
			logrus.SetFormatter(&logrus.TextFormatter{
				ForceColors: true,
			})
			logrus.SetOutput(colorable.NewColorableStdout()) // for windows color output
			if verbose {
				logrus.SetLevel(logrus.DebugLevel)
			}
			return c.run()
		},
	}

	f := cmd.Flags()
	f.BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
	f.StringVarP(&c.configPath, "config", "c", "configs/default_capui_config.yaml", "path of config file (yaml)")
	f.IntVarP(&c.uiPort, "port", "p", 8080, "port of web ui serve port")

	return
}

func (cmd *capUiCmd) run() error {
	if c, err := comm.GetConfig(cmd.configPath); err != nil {
		return err
	} else {
		server := server.NewCapUiServer(c)
		return server.Run(fmt.Sprintf(":%v", cmd.uiPort))
	}
}
