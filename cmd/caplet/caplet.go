package main

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
)

func main() {
	r := gin.Default()

	r.GET("/fetch/:host/:repo", func(c *gin.Context) {
		host := c.Param("host")
		repo := c.Param("repo")

		// docker client
		ctx := context.Background()

		cliVsn, exist := os.LookupEnv("DOCKER_CLIENT_VERSION")
		if !exist {
			cliVsn = "1.39"
		}

		cli, err := client.NewClientWithOpts(client.WithVersion(cliVsn))
		if err != nil {
			c.AbortWithError(500, err)
			return
		}

		// docker pull image
		image := host + "/" + repo
		log.Println("pulling image: ", image)
		out, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
		if err != nil {
			log.Println("pull image failed: ", image)
			c.AbortWithError(500, err)
			return
		}
		defer out.Close()

		// log to console
		termFd, isTerm := term.GetFdInfo(os.Stderr)
		ws := io.MultiWriter(os.Stdout, c.Writer)
		jsonmessage.DisplayJSONMessagesStream(out, ws, termFd, isTerm, nil)
		log.Println("pulled image: ", image)

		// response
		c.String(200, "image: [%s] pull complete", image)
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
