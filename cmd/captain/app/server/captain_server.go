package server

import "io"

type CaptainServer struct {
	out       io.Writer
	endpoints []string
	port      int
}

func NewCaptainServer(out io.Writer, endpoints []string, port int) *CaptainServer {
	return &CaptainServer{
		out:       out,
		endpoints: endpoints,
		port:      port,
	}
}
