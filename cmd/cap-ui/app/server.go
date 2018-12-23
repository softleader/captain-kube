package app

import (
	"github.com/Sirupsen/logrus"
	"github.com/softleader/captain-kube/cmd/cap-ui/app/server"
	"github.com/softleader/captain-kube/cmd/cap-ui/app/server/comm"
	"github.com/spf13/cobra"
)

type capuiCmd struct {
	log        *logrus.Logger
	configPath string
	port       int
}

func NewCapuiCommand() (cmd *cobra.Command) {
	var verbose bool
	c := capuiCmd{}
	cmd = &cobra.Command{
		Use:  "capui",
		Long: "capui is a web interface for captain",
		RunE: func(cmd *cobra.Command, args []string) error {
			c.log = logrus.New()
			c.log.SetOutput(cmd.OutOrStdout())
			if verbose {
				c.log.SetLevel(logrus.DebugLevel)
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

func (cmd *capuiCmd) run() error {
	c, err := comm.GetConfig(cmd.configPath)
	if err != nil {
		return err
	} else {
		return server.Ui(c, cmd.port)
	}
}
