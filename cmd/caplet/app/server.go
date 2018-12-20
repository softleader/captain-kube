package app

import (
	"fmt"
	"github.com/softleader/captain-kube/cmd/caplet/app/server"
	"github.com/softleader/captain-kube/pkg/caplet"
	"github.com/softleader/captain-kube/pkg/env"
	"github.com/softleader/captain-kube/pkg/verbose"
	"github.com/spf13/cobra"
	"io"
)

type capletCmd struct {
	out   io.Writer
	serve string
	port  int
}

type Caplet interface {
	Serve(out io.Writer, port int) error
}

var servers = map[string]Caplet{
	"grpc": server.Grpc{},
	"rest": server.Rest{},
}

func newCapletCmd() (c *capletCmd) {
	c = &capletCmd{
		port:  env.LookupInt(caplet.EnvPort, caplet.DefaultPort),
		serve: env.Lookup(caplet.EnvServe, caplet.DefaultServe),
	}
	return
}

func NewCapletCommand() (cmd *cobra.Command) {
	c := newCapletCmd()
	cmd = &cobra.Command{
		Use:  "caplet",
		Long: "caplet is a daemon run on every kubernetes node",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.run()
		},
	}

	c.out = cmd.OutOrStdout()
	f := cmd.Flags()
	f.BoolVarP(&verbose.Enabled, "verbose", "v", verbose.Enabled, "enable verbose output")
	f.StringVar(&c.serve, "serve", c.serve, "specify api kind to serve (grpc or rest), override "+caplet.EnvServe)
	f.IntVarP(&c.port, "port", "p", c.port, "specify the port to serve, override "+caplet.EnvPort)

	return
}

func (c *capletCmd) run() error {
	if s, err := retrieveServer(c.serve); err != nil {
		return err
	} else {
		return s.Serve(c.out, c.port)
	}
}

func retrieveServer(kind string) (Caplet, error) {
	s, found := servers[kind]
	if !found {
		return nil, fmt.Errorf("unsupported serve kind: %s", kind)
	}
	return s, nil
}
