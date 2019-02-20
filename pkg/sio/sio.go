package sio // stand for streaming io
import "github.com/gin-gonic/gin"

// StreamWriter 定義了 streaming 形式的 writer
type StreamWriter struct {
	send func(p []byte) error
}

// Write 實作 io.writer
func (s *StreamWriter) Write(p []byte) (n int, err error) {
	if err = s.send(p); err == nil {
		n = len(p)
	}
	return
}

// NewStreamWriter 建立一個 StreamWriter
func NewStreamWriter(send func(p []byte) error) *StreamWriter {
	return &StreamWriter{send}
}

// NewSSEventWriter 依照傳入的 gin context 建立 server send stream writer
func NewSSEventWriter(c *gin.Context) *StreamWriter {
	return NewStreamWriter(func(p []byte) error {
		c.SSEvent("message", string(p))
		return nil
	})
}
