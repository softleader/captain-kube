package sse

import (
	"github.com/gin-gonic/gin"
)

func NewWriter(c *gin.Context) *Writer {
	return &Writer{
		GinContext: c,
	}
}

type Writer struct {
	GinContext *gin.Context
}

func (w *Writer) Write(p []byte) (n int, err error) {
	SSE(w.GinContext, string(p))
	return len(p), nil
}

func (w *Writer) WriteStr(s string) (n int, err error) {
	SSE(w.GinContext, s)
	return len(s), nil
}
