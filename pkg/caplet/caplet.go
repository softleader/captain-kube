package caplet

import (
	"context"
	"fmt"
	"github.com/softleader/captain-kube/pkg/dur"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/verbose"
	"google.golang.org/grpc"
	"io"
	"strings"
	"sync"
)

const (
	EnvPort         = "CAPLET_PORT"
	EnvServe        = "CAPLET_SERVE"
	EnvHostname     = "CAPLET_HOSTNAME"
	DefaultPort     = 50051
	DefaultServe    = "grpc"
	DefaultHostname = "caplet"
)

func pullImage(out io.Writer, endpoint string, port int, req *proto.PullImageRequest, timeout int64) error {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%v", endpoint, port), grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("did not connect: %v", err)
	}
	defer conn.Close()
	c := proto.NewCapletClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), dur.Deadline(timeout))
	defer cancel()
	r, err := c.PullImage(ctx, req)
	if err != nil {
		return fmt.Errorf("could not pull image: %v", err)
	}
	if verbose.Enabled {
		for _, i := range r.Images {
			fmt.Fprintf(out, "pulled %v", i)
		}
	}
	return nil
}

func PullImage(out io.Writer, endpoints []string, port int, req *proto.PullImageRequest, timeout int64) error {
	ch := make(chan error, len(endpoints))
	var wg sync.WaitGroup
	for _, ep := range endpoints {
		wg.Add(1)
		go func(out io.Writer, endpoint string, port int, req *proto.PullImageRequest, timeout int64) {
			defer wg.Done()
			ch <- pullImage(out, endpoint, port, req, timeout)
		}(out, ep, port, req, timeout)
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
