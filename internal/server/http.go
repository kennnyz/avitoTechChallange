package server

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func NewHTTPServer(addr string, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

func (s *Server) Run() error {
	logrus.Println("starting server at ", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}
