package caplet

import (
	"context"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/dur"
	"github.com/softleader/captain-kube/pkg/proto"
	"google.golang.org/grpc"
	"io"
)

func (e *Endpoint) PullImage(log *logrus.Logger, req *proto.PullImageRequest, timeout int64) error {
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
		log.Writer().Write(recv.GetMsg())
	}
	return nil
}