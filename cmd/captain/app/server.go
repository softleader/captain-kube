package app

import (
	"fmt"
	"github.com/softleader/captain-kube/cmd/captain/app/server"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/verbose"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
	"net"
)

var ErrNonCapletDaemonFound = fmt.Errorf("non caplet daemon found")

type captainCmd struct {
	out            io.Writer
	serve          string
	endpoints      []string
	port           int
	capletHostname string
	capletPort     int
}

func NewCaptainCommand() (cmd *cobra.Command) {
	c := captainCmd{}
	cmd = &cobra.Command{
		Use:  "captain",
		Long: "captain is the brain of captain-kube system",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(c.endpoints) == 0 {
				if c.endpoints, err = net.LookupHost(c.capletHostname); err != nil {
					return
				}
			}
			if len(c.endpoints) == 0 {
				return ErrNonCapletDaemonFound
			}
			return c.run()
		},
	}

	c.out = cmd.OutOrStdout()
	f := cmd.Flags()
	f.BoolVarP(&verbose.Enabled, "verbose", "v", verbose.Enabled, "enable verbose output")
	f.IntVarP(&c.port, "port", "p", 8081, "specify the port of captain")
	f.StringVar(&c.capletHostname, "caplet-hostname", "caplet", "specify the hostname of caplet daemon to lookup if '--caplet-endpoint' is not specified")
	f.IntVar(&c.capletPort, "caplet-port", 50051, "specify the port of caplet daemon to connect")
	f.StringArrayVar(&c.endpoints, "caplet-endpoint", []string{"localhost"}, "specify the endpoint of caplet daemon to connect")

	return
}

func (c *captainCmd) run() (err error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", c.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterCaptainServer(s, server.NewCaptainServer(c.out, c.endpoints, c.capletPort))
	reflection.Register(s)
	verbose.Fprintf(c.out, "listening and serving GRPC on %v\n", lis.Addr().String())
	return s.Serve(lis)
}
