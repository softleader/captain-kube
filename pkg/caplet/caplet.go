package caplet

import (
	"context"
	"fmt"
	"github.com/softleader/captain-kube/pkg/dur"
	"github.com/softleader/captain-kube/pkg/logger"
	"github.com/softleader/captain-kube/pkg/proto"
	"google.golang.org/grpc"
	"io"
	"strings"
	"sync"
)

const (
	EnvPort         = "CAPLET_PORT"
	EnvHostname     = "CAPLET_HOSTNAME"
	DefaultPort     = 50051
	DefaultHostname = "caplet"
)

type Endpoint struct {
	Target string
	Port   int
}

func (e *Endpoint) PullImage(log *logger.Logger, req *proto.PullImageRequest, timeout int64) error {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%v", e.Target, e.Port), grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("[%s] did not connect: %v", e.Target, err)
	}
	defer conn.Close()
	c := proto.NewCapletClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), dur.Deadline(timeout))
	defer cancel()
	stream, err := c.PullImage(ctx, req)
	if err != nil {
		return fmt.Errorf("[%s] could not pull image: %v", e.Target, err)
	}
	for {
		recv, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Errorf("%v.PullImage(_) = _, %v", c, err)
		}
		log.GetOutput().Write(recv.GetMsg())
	}
	return nil
}

func PullImage(log *logger.Logger, endpoints []*Endpoint, req *proto.PullImageRequest, timeout int64) error {
	ch := make(chan error, len(endpoints))
	var wg sync.WaitGroup
	for _, ep := range endpoints {
		wg.Add(1)
		go func(log *logger.Logger, endpoint *Endpoint, req *proto.PullImageRequest, timeout int64) {
			defer wg.Done()
			ch <- ep.PullImage(log, req, timeout)
		}(log, ep, req, timeout)
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
