package main

import (
	"context"
	"github.com/softleader/captain-kube/pkg/image"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := image.NewPullerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Pull(ctx, &image.PullRequest{
		Host:  "softleader",
		Repo:  "caplet",
	})
	if err != nil {
		log.Fatalf("could not pull image: %v", err)
	}
	log.Printf("Pull %s tag: %s", r.GetResults().GetTag(), r.GetResults().GetMessage())
}
