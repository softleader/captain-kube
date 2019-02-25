package caplet

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	pb "github.com/softleader/captain-kube/pkg/proto"
	"google.golang.org/grpc"
	"io"
	"time"
)

// CallRmi 呼叫 caplet Rmi gRPC api
func (e *Endpoint) CallRmi(log *logrus.Logger, req *pb.RmiRequest, timeout time.Duration) error {
	conn, err := grpc.Dial(e.String(), grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("[%s] did not connect: %v", e.Target, err)
	}
	defer conn.Close()
	c := pb.NewCapletClient(conn)
	log.Debugf("setting context with timeout %v", timeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	stream, err := c.Rmi(ctx, req)
	if err != nil {
		return fmt.Errorf("[%s] %v.Rmi(%v) = _, %v", e.Target, c, req, err)
	}
	var last *pb.ChunkMessage
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
