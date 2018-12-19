package app

import (
	"fmt"
	"github.com/softleader/captain-kube/cmd/caplet/app/server"
	"github.com/spf13/cobra"
	"io"
)

type capletCmd struct {
	out     io.Writer
	verbose bool
	serve   string
	port    int
}

type Caplet interface {
	Serve(out io.Writer, port int) error
}

var servers = map[string]Caplet{
	"grpc": server.Grpc{},
	"rest": server.Rest{},
}

func NewCapletCommand() (cmd *cobra.Command) {
	c := capletCmd{}
	cmd = &cobra.Command{
		Use:  "caplet",
		Long: "caplet is a daemon run on nodes",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.run()
		},
	}

	c.out = cmd.OutOrStdout()
	f := cmd.Flags()
	f.BoolVarP(&c.verbose, "verbose", "v", c.verbose, "enable verbose output, Overrides $SL_VERBOSE")
	f.StringVarP(&c.serve, "serve", "s", "grpc", "determine which kind of api to serve (grpc or rest)")
	f.IntVarP(&c.port, "port", "p", 50051, "determine which port to serve")

	return
}

func (c *capletCmd) run() error {
	s, found := servers[c.serve]
	if !found {
		return fmt.Errorf("unsupported serve type: %s", c.serve)
	}
	return s.Serve(c.out, c.port)
}
