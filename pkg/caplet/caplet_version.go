package caplet

import (
	"context"
	"fmt"
	"github.com/softleader/captain-kube/pkg/proto"
	"google.golang.org/grpc"
	"time"
)

func (e *Endpoint) Version(full bool, timeout time.Duration) (*captainkube_v2.VersionResponse, error) {
	conn, err := grpc.Dial(e.String(), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("[%s] did not connect: %v", e.Target, err)
	}
	defer conn.Close()
	c := captainkube_v2.NewCapletClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return c.Version(ctx, &captainkube_v2.VersionRequest{
		Full: full,
	})
}
