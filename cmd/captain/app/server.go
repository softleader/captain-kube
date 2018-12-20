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

type CaptainCmd struct {
	Out            io.Writer
	Serve          string
	Endpoints      []string
	Port           int
	CapletHostname string
	CapletPort     int
}

func NewCaptainCommand() (cmd *cobra.Command) {
	c := CaptainCmd{
		Port:           env.LookupInt(captain.EnvPort, captain.DefaultPort),
		CapletPort:     env.LookupInt(caplet.EnvPort, caplet.DefaultPort),
		CapletHostname: env.Lookup(caplet.EnvHostname, caplet.DefaultHostname),
	}
	cmd = &cobra.Command{
		Use:  "captain",
		Long: "captain is the brain of captain-kube system",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			return c.Run()
		},
	}

	c.Out = cmd.OutOrStdout()
	f := cmd.Flags()
	f.BoolVarP(&verbose.Enabled, "verbose", "v", verbose.Enabled, "enable verbose output")
	f.IntVarP(&c.Port, "port", "p", c.Port, "specify the port of captain, override "+captain.EnvPort)
	f.StringVar(&c.CapletHostname, "caplet-hostname", c.CapletHostname, "specify the hostname of caplet daemon to lookup, override "+caplet.EnvHostname)
	f.IntVar(&c.CapletPort, "caplet-port", c.CapletPort, "specify the port of caplet daemon to connect, override "+caplet.EnvPort)
	f.StringArrayVarP(&c.Endpoints, "caplet-endpoint", "e", []string{""}, "specify the endpoint of caplet daemon to connect, override '--caplet-hostname'")

	return
}

func (c *CaptainCmd) Run() (err error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", c.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterCaptainServer(s, server.NewCaptainServer(c.Out, c.CapletHostname, c.Endpoints, c.CapletPort))
	reflection.Register(s)
	verbose.Fprintf(c.Out, "listening and serving GRPC on %v\n", lis.Addr().String())
	return s.Serve(lis)
}
