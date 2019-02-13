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

func (e *Endpoint) Rmi(log *logrus.Logger, req *captainkube_v2.RmiRequest, timeout int64) error {
	conn, err := grpc.Dial(e.String(), grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("[%s] did not connect: %v", e.Target, err)
	}
	defer conn.Close()
	c := captainkube_v2.NewCapletClient(conn)
	deadline := dur.Deadline(timeout)
	log.Debugf("setting context with timeout %v", deadline)
	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()
	stream, err := c.Rmi(ctx, req)
	if err != nil {
		return fmt.Errorf("[%s] %v.Rmi(%v) = _, %v", e.Target, c, req, err)
	}
	var last *captainkube_v2.ChunkMessage
	for {
		recv, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("[%s] failed to receive a chunk msg : %v", e.Target, err)
		}
		log.Out.Write(e.Color(format(last, recv)))
		last = recv
	}
	return nil
}
