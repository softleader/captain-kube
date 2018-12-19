package server

import (
	"fmt"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
	"github.com/gin-gonic/gin"
	"github.com/softleader/captain-kube/cmd/caplet/app/dockerctl"
	"io"
	"log"
	"os"
	"strings"
)

type Rest struct{}

func (_ Rest) Serve(out io.Writer, port int) error {
	r := gin.Default()

	r.GET("/fetch/:host/:repotag", func(c *gin.Context) {
		host := c.Param("host")
		repotag := c.Param("repotag")
		image := host + repotag

		// split report tag
		ss := strings.Split(repotag, ":")
		repo := ss[0]
		var tag string
		if len(ss) >= 2 {
			tag = ss[1]
		} else {
			tag = "latest"
		}

		// pull image
		out, err := dockerctl.Pull(host, repo, tag)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		defer out.Close()

		// log to console and response
		termFd, isTerm := term.GetFdInfo(os.Stderr)
		ws := io.MultiWriter(os.Stdout, c.Writer)
		jsonmessage.DisplayJSONMessagesStream(out, ws, termFd, isTerm, nil)
		log.Println("pulled image: ", image)

		// response
		c.String(200, "image: [%s] pull complete", image)
	})

	return r.Run(fmt.Sprintf(":%v", port))
}
