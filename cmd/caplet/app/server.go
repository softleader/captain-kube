package app

import (
	"fmt"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/cmd/caplet/app/server"
	"github.com/softleader/captain-kube/pkg/caplet"
	"github.com/softleader/captain-kube/pkg/env"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/release"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type capletCmd struct {
	metadata *release.Metadata
	serve    string
	port     int
}

func NewCapletCommand(metadata *release.Metadata) (cmd *cobra.Command) {
	var verbose bool
	c := &capletCmd{
		metadata: metadata,
		port:     env.LookupInt(caplet.EnvPort, caplet.DefaultPort),
	}
	cmd = &cobra.Command{
		Use:  "caplet",
		Long: "caplet is a daemon run on every kubernetes node",
		RunE: func(cmd *cobra.Command, args []string) error {
			logrus.SetOutput(colorable.NewColorableStdout()) // for windows color output
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
	captainkube_v2.RegisterCapletServer(s, server.NewCapletServer(c.metadata))
	logrus.Printf("registered caplet server\n")
	reflection.Register(s)

	logrus.Printf("listening and serving GRPC on %v\n", lis.Addr().String())
	return s.Serve(lis)
}
