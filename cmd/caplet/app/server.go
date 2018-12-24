package app

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/softleader/captain-kube/cmd/caplet/app/server"
	"github.com/softleader/captain-kube/pkg/caplet"
	"github.com/softleader/captain-kube/pkg/env"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type capletCmd struct {
	serve string
	port  int
}

func newCapletCmd() (c *capletCmd) {
	c = &capletCmd{
		port: env.LookupInt(caplet.EnvPort, caplet.DefaultPort),
	}
	return
}

func NewCapletCommand() (cmd *cobra.Command) {
	var verbose bool
	c := newCapletCmd()
	cmd = &cobra.Command{
		Use:  "caplet",
		Long: "caplet is a daemon run on every kubernetes node",
		RunE: func(cmd *cobra.Command, args []string) error {
			logrus.SetOutput(cmd.OutOrStdout())
			logrus.SetFormatter(&logrus.TextFormatter{
				ForceColors: true,
			})
			if verbose {
				logrus.SetLevel(logrus.DebugLevel)
			}
			return c.run()
		},
	}

	f := cmd.Flags()
	f.BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
	f.IntVarP(&c.port, "port", "p", c.port, "specify the port to serve, override "+caplet.EnvPort)

	return
}

func (c *capletCmd) run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", c.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterCapletServer(s, server.NewCapletServer())
	logrus.Printf("registered caplet server\n")
	reflection.Register(s)

	logrus.Printf("listening and serving GRPC on %v\n", lis.Addr().String())
	return s.Serve(lis)
}
