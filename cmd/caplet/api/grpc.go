package api

import (
	"context"
	"github.com/softleader/captain-kube/cmd/caplet/dockerctl"
	"github.com/softleader/captain-kube/pkg/image"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io/ioutil"
	"log"
	"net"
)

const (
	port = ":50051"
)

type server struct{}

func (s *server) Pull(_ context.Context, req *image.PullRequest) (*image.PullResponse, error) {

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

func Grpc() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	image.RegisterPullerServer(s, &server{})

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
