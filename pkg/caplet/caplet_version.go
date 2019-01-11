package caplet

import (
	"context"
	"fmt"
	"github.com/softleader/captain-kube/pkg/dur"
	"github.com/softleader/captain-kube/pkg/proto"
	"google.golang.org/grpc"
)

func (e *Endpoint) Version(full bool, timeout int64) (*captainkube_v2.VersionResponse, error) {
	conn, err := grpc.Dial(e.String(), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("[%s] did not connect: %v", e.Target, err)
	}
	defer conn.Close()
	c := captainkube_v2.NewCapletClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), dur.Deadline(timeout))
	defer cancel()
	return c.Version(ctx, &captainkube_v2.VersionRequest{
		Full: full,
	})
}
