package captain

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/proto"
	"google.golang.org/grpc"
	"time"
)

// CallConsoleURL 呼叫 captain ConsoleURL gRPC API
func CallConsoleURL(log *logrus.Logger, url string, req *captainkube_v2.ConsoleURLRequest, timeout time.Duration) (*captainkube_v2.ConsoleURLResponse, error) {
	log.Debugf("dialing %q with insecure", url)
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("did not connect: %v", err)
	}
	defer conn.Close()
	c := captainkube_v2.NewCaptainClient(conn)
	log.Debugf("setting context with timeout %v", timeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resq, err := c.ConsoleURL(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%v.ConsoleURL(%v) = _, %v", c, req, err)
	}
	return resq, nil
}
