package pas

import (
	"log"
	"net/http"

	"github.com/rs/cors"
)

type Server struct {
	http.Server
}

func NewServer(addr string, handler http.Handler) *Server {
	return &Server{
		Server: http.Server{
			Addr:    addr,
			Handler: cors.Default().Handler(handler),
		},
	}
}

func (s *Server) ListenAndServe() {
	err := s.Server.ListenAndServe()
	if err == http.ErrServerClosed {
		return
	}
	log.Fatal(err)
}
