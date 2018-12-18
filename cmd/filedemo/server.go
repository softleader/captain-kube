package main

import (
	"fmt"
	"github.com/softleader/captain-kube/pkg/proto/file"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
	"log"
	"net"
	"os"
)

const (
	port = ":50052"
)

type server struct{}

func (s *server) Upload(stream file.Uploader_UploadServer) error {
	log.Println("reciving file")

	os.Remove("./data/uploaded.txt")
	f, err := os.Create("./data/uploaded.txt")
	if err != nil {
		return fmt.Errorf("failed to create file: %s", err)
	}

	for {
		c, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return fmt.Errorf("failed unexpectadely while reading chunks from stream: %s", err)
			}
		}
		b := c.GetContent()
		log.Println("reciving", len(b), "byte")
		f.Write(b)
	}

	info, _ := f.Stat()
	log.Println("file recived: ", info.Size(), "bytes")
	defer f.Close()

	return stream.SendAndClose(&file.UploadResponse{
		Message: "Upload received with success",
	})
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	file.RegisterUploaderServer(s, &server{})

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
