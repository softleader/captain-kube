package captain

import (
	"context"
	"fmt"
	"github.com/softleader/captain-kube/pkg/dur"
	"github.com/softleader/captain-kube/pkg/proto"
	"google.golang.org/grpc"
	"io"
)

const (
	EnvPort     = "CAPTAIN_PORT"
	DefaultPort = 8081
)

func InstallChart(out io.Writer, url string, req *proto.InstallChartRequest, timeout int64) error {

	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("did not connect: %v", err)
	}
	defer conn.Close()
	c := proto.NewCaptainClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), dur.Deadline(timeout))
	defer cancel()
	r, err := c.InstallChart(ctx, req)
	if err != nil {
		return fmt.Errorf("could not install chart: %v", err)
	}
	fmt.Fprintf(out, "chart installed %v", r.GetOut())
	return nil

}

func GenerateScript(out io.Writer, url string, req *proto.GenerateScriptRequest, timeout int64) error {

	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("did not connect: %v", err)
	}
	defer conn.Close()
	c := proto.NewCaptainClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), dur.Deadline(timeout))
	defer cancel()
	r, err := c.GenerateScript(ctx, req)
	if err != nil {
		return fmt.Errorf("could not generate script: %v", err)
	}
	fmt.Fprintf(out, "script generated %v", r.GetOut())
	return nil

}
