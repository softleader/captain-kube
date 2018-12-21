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
type server struct {
	out io.Writer
}

func (g *server) PullImage(req *proto.PullImageRequest, stream proto.Caplet_PullImageServer) error {
	sout := sio.NewStreamWriter(func(p []byte) error {
		return stream.Send(&proto.PullImageResponse{
			Msg: p,
		})
	})
	for _, image := range req.Images {
		if err := pull(g.out, sout, image, req.GetRegistryAuth()); err != nil {
			return err
		}
	}
	return nil
}

func pull(out io.Writer, sout io.Writer, image *proto.Image, auth *proto.RegistryAuth) error {
	if tag := image.GetTag(); len(tag) == 0 {
		image.Tag = "latest"
	}
	rc, err := dockerctl.Pull(out, chart.Image{
		Host: image.Host,
		Repo: image.Repo,
		Tag:  image.Tag,
	}, auth)
	if err != nil {
		return err
	}
	defer rc.Close()
	return jsonmessage.DisplayJSONMessagesToStream(rc, command.NewOutStream(sout), nil)
}

func (_ Grpc) Serve(out io.Writer, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterCapletServer(s, &server{out})
	verbose.Fprintf(out, "registered caplet server\n")
	reflection.Register(s)

	verbose.Fprintf(out, "listening and serving GRPC on %v\n", lis.Addr().String())
	return s.Serve(lis)
}
