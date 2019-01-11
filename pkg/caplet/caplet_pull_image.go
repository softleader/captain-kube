package caplet

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/dur"
	"github.com/softleader/captain-kube/pkg/proto"
	"google.golang.org/grpc"
	"io"
)

func (e *Endpoint) PullImage(log *logrus.Logger, req *captainkube_v2.PullImageRequest, timeout int64) error {
	conn, err := grpc.Dial(e.String(), grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("[%s] did not connect: %v", e.Target, err)
	}
	defer conn.Close()
	c := captainkube_v2.NewCapletClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), dur.Deadline(timeout))
	defer cancel()
	stream, err := c.PullImage(ctx, req)
	if err != nil {
		return fmt.Errorf("[%s] could not pull image: %v", e.Target, err)
	}
	var last *captainkube_v2.ChunkMessage
	for {
		recv, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("%v.PullImage(_) = _, %v", c, err)
		}
		log.Out.Write(e.Color(format(last, recv)))
		last = recv
	}
	return nil
}
