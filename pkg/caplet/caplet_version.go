package caplet

import (
	"context"
	"fmt"
	pb "github.com/softleader/captain-kube/pkg/proto"
	"google.golang.org/grpc"
	"time"
)

// CallVersion 呼叫 caplet Version gRPC api
func (e *Endpoint) CallVersion(full bool, timeout time.Duration) (*pb.VersionResponse, error) {
	conn, err := grpc.Dial(e.String(), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("[%s] did not connect: %v", e.Target, err)
	}
	defer conn.Close()
	c := pb.NewCapletClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return c.Version(ctx, &pb.VersionRequest{
		Full: full,
	})
}
