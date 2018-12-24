package sio // stand for streaming io

type StreamWriter struct {
	send func(p []byte) error
}

func (s *StreamWriter) Write(p []byte) (n int, err error) {
	if err = s.send(p); err == nil {
		n = len(p)
	}
	return
}

func NewStreamWriter(send func(p []byte) error) *StreamWriter {
	return &StreamWriter{send}
}