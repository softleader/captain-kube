package server

import (
	"context"
	"fmt"
	"github.com/softleader/captain-kube/cmd/caplet/app/dockerctl"
	"github.com/softleader/captain-kube/pkg/proto/image"
	"github.com/softleader/captain-kube/pkg/verbose"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
	"io/ioutil"
	"net"
)

type Grpc struct{}
type server struct{}

func (g *server) Pull(_ context.Context, req *image.PullRequest) (*image.PullResponse, error) {

	var tag string
	if tag = req.GetTag(); len(tag) == 0 {
		tag = "latest"
	}

	out, err := dockerctl.Pull(req.GetHost(), req.GetRepo(), tag)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(out)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	return &image.PullResponse{
		Results: &image.Result{
			Tag:     tag,
			Message: string(bytes),
		},
	}, nil
}

func (_ Grpc) Serve(out io.Writer, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	image.RegisterPullerServer(s, &server{})
	verbose.Fprintf(out, "registered %q\n", "puller service")

	reflection.Register(s)

	verbose.Fprintf(out, "listening on %v\n", lis.Addr().String())
	return s.Serve(lis)
}
