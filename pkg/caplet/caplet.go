package caplet

import (
	"context"
	"fmt"
	"github.com/softleader/captain-kube/pkg/dur"
	"github.com/softleader/captain-kube/pkg/proto"
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

type Endpoint struct {
	Target  string
	Port    int
	Timeout int64
}

func (e *Endpoint) PullImage(out io.Writer, req *proto.PullImageRequest) error {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%v", e.Target, e.Port), grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("did not connect: %v", err)
	}
	defer conn.Close()
	c := proto.NewCapletClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), dur.Deadline(e.Timeout))
	defer cancel()
	stream, err := c.PullImage(ctx, req)
	if err != nil {
		return fmt.Errorf("could not pull image: %v", err)
	}
	for {
		recv, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Errorf("%v.PullImage(_) = _, %v", c, err)
		}
		out.Write(recv.GetMessage())
	}
	return nil
}

func PullImage(out io.Writer, endpoints []*Endpoint, req *proto.PullImageRequest) error {
	ch := make(chan error, len(endpoints))
	var wg sync.WaitGroup
	for _, ep := range endpoints {
		wg.Add(1)
		go func(out io.Writer, endpoint *Endpoint, req *proto.PullImageRequest) {
			defer wg.Done()
			ch <- ep.PullImage(out, req)
		}(out, ep, req)
	}
	wg.Wait()
	close(ch)
	var errors []string
	for e := range ch {
		if e != nil {
			errors = append(errors, e.Error())
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "\n"))
	}
	return nil
}
