package server

import (
	"fmt"
	"github.com/docker/docker/cli/command"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/softleader/captain-kube/pkg/dockerctl"
	"github.com/softleader/captain-kube/pkg/helm/chart"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/sio"
	"github.com/softleader/captain-kube/pkg/verbose"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
	"net"
)

type Grpc struct{}
type server struct{}

func (g *server) PullImage(req *proto.PullImageRequest, stream proto.Caplet_PullImageServer) error {
	for _, image := range req.Images {
		sw := sio.NewStreamWriter(func(p []byte) error {
			return stream.Send(&proto.PullImageResponse{
				Msg: p,
			})
		})
		if err := pull(sw, image, req.GetRegistryAuth()); err != nil {
			return err
		}
	}
	return nil
}

func pull(w io.Writer, image *proto.Image, auth *proto.RegistryAuth) error {
	if tag := image.GetTag(); len(tag) == 0 {
		image.Tag = "latest"
	}
	out, err := dockerctl.Pull(w, chart.Image{
		Host: image.Host,
		Repo: image.Repo,
		Tag:  image.Tag,
	}, auth)
	if err != nil {
		return err
	}
	defer out.Close()
	return jsonmessage.DisplayJSONMessagesToStream(out, command.NewOutStream(sw), nil)
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
