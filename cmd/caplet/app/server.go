package app

import (
	"fmt"
	"github.com/softleader/captain-kube/cmd/caplet/app/server"
	"github.com/softleader/captain-kube/pkg/verbose"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strconv"
)

const (
	EnvCapletPort  = "CAPLET_PORT"
	EnvCapletServe = "CAPLET_SERVE"

	defaultEnvCapletPort  = 50051
	defaultEnvCapletServe = "grpc"
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

func NewCapletCommand() (cmd *cobra.Command) {
	c := capletCmd{
		port:  defaultEnvCapletPort,
		serve: defaultEnvCapletServe,
	}
	if value, found := os.LookupEnv(EnvCapletPort); found {
		if p, err := strconv.Atoi(value); err != nil {
			c.port = p
		}
	}
	if value, found := os.LookupEnv(EnvCapletServe); found {
		c.serve = value
	}

	cmd = &cobra.Command{
		Use:  "caplet",
		Long: "caplet is a daemon run on nodes",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.run()
		},
	}

	c.out = cmd.OutOrStdout()
	f := cmd.Flags()
	f.BoolVarP(&verbose.Verbose, "verbose", "v", verbose.Verbose, "enable verbose output")
	f.StringVar(&c.serve, "serve", c.serve, "specify api kind to serve (grpc or rest), override "+EnvCapletServe)
	f.IntVarP(&c.port, "port", "p", c.port, "specify the port to serve, override "+EnvCapletPort)

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
