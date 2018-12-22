package sse

import (
	"github.com/gin-gonic/gin"
)

func NewWriter(c *gin.Context) *SSEWriter {
	return &SSEWriter{
		GinContext: c,
	}
}

type SSEWriter struct {
	GinContext *gin.Context
}

func (w *SSEWriter) Write(p []byte) (n int, err error) {
	SSE(w.GinContext, string(p))
	return len(p), nil
}

func (w *SSEWriter) WriteStr(s string) (n int, err error) {
	SSE(w.GinContext, s)
	return len(s), nil
}
