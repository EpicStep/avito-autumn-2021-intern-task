package server

import (
	"context"
	"net/http"
)

// Server is a http server.
type Server struct {
	server *http.Server
}

// New returns a new http server.
func New(addr string, router http.Handler) *Server {
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return &Server{
		server: srv,
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
