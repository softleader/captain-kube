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

func GenerateScript(log *logrus.Logger, url string, req *tw_com_softleader.GenerateScriptRequest, timeout int64) error {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("did not connect: %v\n", err)
	}
	defer conn.Close()
	c := tw_com_softleader.NewCaptainClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), dur.Deadline(timeout))
	defer cancel()
	stream, err := c.GenerateScript(ctx, req)
	if err != nil {
		return fmt.Errorf("could not generate script: %v\n", err)
	}
	for {
		recv, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("%v.GenerateScript(_) = _, %v", c, err)
		}
		log.Out.Write(recv.GetMsg())
	}
	return nil

}
