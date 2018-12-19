package caplet

import (
	"context"
	"fmt"
	"github.com/softleader/captain-kube/pkg/dur"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/verbose"
	"google.golang.org/grpc"
	"io"
)

func PullImage(out io.Writer, endpoint string, port int, req *proto.PullImageRequest, timeout int) error {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%v", endpoint, port), grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("did not connect: %v", err)
	}
	defer conn.Close()
	c := proto.NewCapletClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), dur.Deadline(timeout))
	defer cancel()
	r, err := c.PullImage(ctx, req)
	if err != nil {
		return fmt.Errorf("could not pull image: %v", err)
	}
	verbose.Fprintf(out, "Pull %s tag:\n%s", r.GetTag(), r.GetMessage())
	return nil
}
