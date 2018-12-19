package app

import (
	"fmt"
	"github.com/softleader/captain-kube/cmd/captain/app/client"
	"github.com/softleader/captain-kube/pkg/verbose"
	"github.com/spf13/cobra"
	"io"
	"net"
	"strings"
	"sync"
)

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
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.run()
		},
	}

	c.out = cmd.OutOrStdout()
	f := cmd.Flags()
	f.BoolVarP(&verbose.Verbose, "verbose", "v", verbose.Verbose, "enable verbose output")
	f.StringArrayVarP(&c.endpoints, "endpoint", "e", []string{}, "specify the endpoints of caplet")
	f.IntVarP(&c.port, "port", "p", 50051, "specify the port of caplet")
	f.StringVar(&c.caplet, "caplet", "caplet", "specify the hostname of caplet daemon to lookup")

	return
}

func (c *captainCmd) run() (err error) {

	if len(c.endpoints) == 0 {
		if c.endpoints, err = net.LookupHost(c.caplet); err != nil {
			return err
		}
	}
	if len(c.endpoints) == 0 {
		return fmt.Errorf("non caplet daemon found")
	}
	ch := make(chan error, len(c.endpoints))
	var wg sync.WaitGroup
	for _, ep := range c.endpoints {
		wg.Add(1)
		go func(out io.Writer, endpoint string, port int) {
			defer wg.Done()
			ch <- client.PullImage(out, endpoint, port)
		}(c.out, ep, c.port)
	}
	wg.Wait()
	close(ch)
	if len(ch) > 0 {
		var errors []string
		for e := range ch {
			errors = append(errors, e.Error())
		}
		return fmt.Errorf(strings.Join(errors, "\n"))
	}
	return nil

}
