package server

import (
	"context"
	"fmt"
	"github.com/softleader/captain-kube/cmd/caplet/app/dockerctl"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/verbose"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
	"io/ioutil"
	"net"
)

type Grpc struct{}
type server struct{}

func (g *server) PullImage(_ context.Context, req *proto.PullImageRequest) (*proto.PullImageResponse, error) {

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

	return &proto.PullImageResponse{
		Tag:     tag,
		Message: string(bytes),
	}, nil
}

func (_ Grpc) Serve(out io.Writer, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterCapletServer(s, &server{})
	verbose.Fprintf(out, "registered caplet server\n")
	reflection.Register(s)

	verbose.Fprintf(out, "listening and serving GRPC on %v\n", lis.Addr().String())
	return s.Serve(lis)
}
