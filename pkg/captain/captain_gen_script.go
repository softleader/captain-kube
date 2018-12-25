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

func GenerateScript(log *logrus.Logger, url string, req *proto.GenerateScriptRequest, timeout int64) error {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("did not connect: %v\n", err)
	}
	defer conn.Close()
	c := proto.NewCaptainClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), dur.Deadline(timeout))
	defer cancel()
	stream, err := c.GenerateScript(ctx, req)
	if err != nil {
		return fmt.Errorf("could not generate script: %v\n", err)
	}
	for {
		recv, err := stream.Recv()
		if err == io.EOF {
			log.Println("### EOF ###")
			break
		}
		if err != nil {
			return fmt.Errorf("%v.GenerateScript(_) = _, %v", c, err)
		}
		log.Writer().Write(recv.GetMsg())
	}
	return nil

}
