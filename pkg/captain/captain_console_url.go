package captain

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/captain-kube/pkg/dur"
	"github.com/softleader/captain-kube/pkg/proto"
	"google.golang.org/grpc"
)

func ConsoleURL(log *logrus.Logger, url string, req *captainkube_v2.ConsoleURLRequest, timeout int64) (*captainkube_v2.ConsoleURLResponse, error) {
	log.Debugf("dialing %q with insecure", url)
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("did not connect: %v", err)
	}
	defer conn.Close()
	c := captainkube_v2.NewCaptainClient(conn)
	deadline := dur.Deadline(timeout)
	log.Debugf("setting context with timeout %v", deadline)
	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()
	resq, err := c.ConsoleURL(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%v.ConsoleURL(%v) = _, %v", c, req, err)
	}
	return resq, nil
}
