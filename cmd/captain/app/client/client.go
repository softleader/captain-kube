package client

import (
	"context"
	"fmt"
	"github.com/softleader/captain-kube/pkg/proto/image"
	"github.com/softleader/slctl/pkg/verbose"
	"google.golang.org/grpc"
	"io"
)

func PullImage(out io.Writer, endpoint string, port int) error {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%v", endpoint, port), grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("did not connect: %v", err)
	}
	defer conn.Close()

	c := image.NewPullerClient(conn)
	ctx := context.Background()
	r, err := c.Pull(ctx, &image.PullRequest{
		Host: "softleader",
		Repo: "caplet",
	})
	if err != nil {
		return fmt.Errorf("could not pull image: %v", err)
	}
	verbose.Fprintf(out, "Pull %s tag:\n%s", r.GetResults().GetTag(), r.GetResults().GetMessage())
	return nil
}
