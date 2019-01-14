package captain

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/dur"
	"github.com/softleader/captain-kube/pkg/proto"
	"google.golang.org/grpc"
	"io"
)

func GenerateScript(log *logrus.Logger, url string, req *captainkube_v2.GenerateScriptRequest, timeout int64) error {
	log.Debugf("dialing %q with insecure", url)
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("did not connect: %v", err)
	}
	defer conn.Close()
	c := captainkube_v2.NewCaptainClient(conn)
	deadline := dur.Deadline(timeout)
	log.Debugf("setting context with timeout %v", deadline)
	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()
	stream, err := c.GenerateScript(ctx, req)
	if err != nil {
		return fmt.Errorf("%v.GenerateScript(%v) = _, %v", c, req, err)
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
