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
	out       io.Writer
	serve     string
	endpoints []string
	port      int
	caplet    string
}

func NewCaptainCommand() (cmd *cobra.Command) {
	c := captainCmd{}
	cmd = &cobra.Command{
		Use:  "captain",
		Long: "captain is the brain of captain-kube system",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(c.endpoints) == 0 {
				if c.endpoints, err = net.LookupHost(c.caplet); err != nil {
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
	f.StringArrayVarP(&c.endpoints, "endpoint", "e", []string{}, "specify the endpoint of caplets")
	f.IntVarP(&c.port, "port", "p", 50051, "specify the port of caplet")
	f.StringVar(&c.caplet, "caplet", "caplet", "specify the hostname of caplet daemon to lookup")

	return
}

func (c *captainCmd) run() (err error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", c.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterCaptainServer(s, server.NewCaptainServer(c.out, c.endpoints, c.port))
	reflection.Register(s)
	verbose.Fprintf(c.out, "listening and serving GRPC on %v\n", lis.Addr().String())
	return s.Serve(lis)
}
