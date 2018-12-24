package server

import (
	"github.com/Sirupsen/logrus"
)

type capletServer struct {
	log *logrus.Logger
}

func NewCapletServer(log *logrus.Logger) (s *capletServer) {
	s = &capletServer{
		log: log,
	}
	return
}
