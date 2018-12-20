package server

import (
	"bytes"
	"context"
	"fmt"
	"github.com/softleader/captain-kube/pkg/dockerctl"
	"github.com/softleader/captain-kube/pkg/proto"
	"github.com/softleader/captain-kube/pkg/verbose"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
	"io/ioutil"
	"net"
	"strings"
)

type Grpc struct{}
type server struct{}

func (g *server) PullImage(_ context.Context, req *proto.PullImageRequest) (*proto.PullImageResponse, error) {

	var buf bytes.Buffer
	var errors []string
	for _, image := range req.Images {
		if out, err := pull(image); err != nil {
			errors = append(errors, err.Error())
		} else {
			if bytes, err := ioutil.ReadAll(out); err != nil { // FIXME: to streaming output
				errors = append(errors, err.Error())
			} else {
				buf.Write(bytes)
				out.Close()
			}
		}
	}

	if len(errors) > 0 {
		return nil, fmt.Errorf(strings.Join(errors, "\n"))
	}

	return &proto.PullImageResponse{
		Images:  req.Images,
		Message: buf.String(),
	}, nil
}
func pull(image *proto.Image) (io.ReadCloser, error) {
	if tag := image.GetTag(); len(tag) == 0 {
		image.Tag = "latest"
	}
	return dockerctl.Pull(image.GetHost(), image.GetRepo(), image.GetTag())
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
