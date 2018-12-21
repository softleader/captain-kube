package sio // stand for streaming io

type StreamWriter struct {
	sender func(p []byte) error
}

func (s *StreamWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	err = s.sender(p)
	return
}

func NewStreamWriter(sender func(p []byte) error) *StreamWriter {
	return &StreamWriter{sender}
}
