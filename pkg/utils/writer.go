package utils

import (
	"github.com/gin-gonic/gin"
)

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
