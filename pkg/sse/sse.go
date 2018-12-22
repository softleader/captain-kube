package sse

import (
	"github.com/gin-gonic/gin"
	"io"
)

func SSE(c *gin.Context, message string) {
	c.Stream(func(w io.Writer) bool {
		c.SSEvent("message", message)
		return false
	})
}
