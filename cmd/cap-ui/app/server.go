package app

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/mattn/go-colorable"
	"github.com/softleader/captain-kube/cmd/cap-ui/app/server"
	"github.com/softleader/captain-kube/cmd/cap-ui/app/server/comm"
	"github.com/spf13/cobra"
)

type capUiCmd struct {
	configPath string
	port       int
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
	f.IntVarP(&c.port, "port", "p", 8080, "port of web ui serve port")

	return
}

func (cmd *capUiCmd) run() error {
	if c, err := comm.GetConfig(cmd.configPath); err != nil {
		return err
	} else {
		server := server.NewCapUiServer(c)
		return server.Run(fmt.Sprintf(":%v", cmd.port))
	}
}
