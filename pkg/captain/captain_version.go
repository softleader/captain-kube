package captain

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	pb "github.com/softleader/captain-kube/pkg/proto"
	"google.golang.org/grpc"
	"io"
	"time"
)

// CallVersion 呼叫 captain Version gRPC API
func CallVersion(log *logrus.Logger, url string, full, color bool, timeout time.Duration) error {
	log.Debugf("dialing %q with insecure", url)
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewCaptainClient(conn)
	log.Debugf("setting context with timeout %v", timeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req := &pb.VersionRequest{
		Full:    full,
		Color:   color,
		Timeout: timeout.String(),
	}
	stream, err := c.Version(ctx, req)
	if err != nil {
		return fmt.Errorf("%v.Version(%v) = _, %v", c, req, err)
	}
	for {
		recv, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to receive a chunk msg : %v", err)
		}
		log.Out.Write(recv.GetMsg())
	}
	return nil

}
