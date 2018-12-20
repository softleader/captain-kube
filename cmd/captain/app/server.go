package app

import (
	"fmt"
	"github.com/softleader/captain-kube/cmd/captain/app/server"
	"github.com/softleader/captain-kube/pkg/caplet"
	"github.com/softleader/captain-kube/pkg/captain"
	"github.com/softleader/captain-kube/pkg/env"
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
	c := captainCmd{
		port:           env.LookupInt(captain.EnvPort, captain.DefaultPort),
		capletPort:     env.LookupInt(caplet.EnvPort, caplet.DefaultPort),
		capletHostname: env.Lookup(caplet.EnvHostname, caplet.DefaultHostname),
	}
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
	f.IntVarP(&c.port, "port", "p", c.port, "specify the port of captain, override "+captain.EnvPort)
	f.StringVar(&c.capletHostname, "caplet-hostname", c.capletHostname, "specify the hostname of caplet daemon to lookup, override "+caplet.EnvHostname)
	f.IntVar(&c.capletPort, "caplet-port", c.capletPort, "specify the port of caplet daemon to connect, override "+caplet.EnvPort)
	f.StringArrayVarP(&c.endpoints, "caplet-endpoint", "e", []string{""}, "specify the endpoint of caplet daemon to connect, override '--caplet-hostname'")

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
