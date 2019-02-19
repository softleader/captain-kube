package caplet

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/proto"
	"google.golang.org/grpc"
	"io"
	"time"
)

func (e *Endpoint) Prune(log *logrus.Logger, req *captainkube_v2.PruneRequest, timeout time.Duration) error {
	conn, err := grpc.Dial(e.String(), grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("[%s] did not connect: %v", e.Target, err)
	}
	defer conn.Close()
	c := captainkube_v2.NewCapletClient(conn)
	log.Debugf("setting context with timeout %v", timeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	stream, err := c.Prune(ctx, req)
	if err != nil {
		return fmt.Errorf("[%s] %v.PullImage(%v) = _, %v", e.Target, c, req, err)
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
