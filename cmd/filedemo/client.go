package main

import (
	"context"
	"github.com/softleader/captain-kube/pkg/proto/file"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
)

func main() {
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := file.NewUploaderClient(conn)
	ctx := context.Background()

	f, err := os.Open("./data/APPLICATION_2.pdf")
	if err != nil {
		log.Fatalf("read file failed: %v", err)
	}

	stream, err := c.Upload(ctx)
	if err != nil {
		log.Fatalf("open upload stream failed: %v", err)
	}

	buf := make([]byte, 4)
	for {
		if n, err := f.Read(buf); err == nil {
			stream.Send(&file.Chunk{
				Content: buf[:n],
			})
		} else {
			if err == io.EOF {
				break
			} else {
				log.Fatalf("uploading stream failed: %v", err)
			}
		}
	}
	stream.CloseAndRecv()

	log.Printf("file uploaded")
}
