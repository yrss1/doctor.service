package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

type Server struct {
	http     *http.Server
	listener net.Listener
}

type Configuration func(s *Server) error

func New(configs ...Configuration) (r *Server, err error) {
	r = &Server{}
	for _, cfg := range configs {
		if err = cfg(r); err != nil {
			return
		}
	}
	return
}

func (s *Server) Run() error {
	if s.http != nil {
		errChan := make(chan error, 1)
		go func() {
			errChan <- s.http.ListenAndServe()
		}()

		select {
		case err := <-errChan:
			if err != http.ErrServerClosed {
				fmt.Printf("ERR_SERVE_HTTP: %v\n", err)
				return err
			}
		}
	}
	return nil
}
func (s *Server) Stop(ctx context.Context) (err error) {
	if s.http != nil {
		if err = s.http.Shutdown(ctx); err != nil {
			return
		}
	}

	return
}
func WithHTTPServer(handler http.Handler, port string) Configuration {
	return func(s *Server) error {
		s.http = &http.Server{
			Addr:    ":" + port,
			Handler: handler,
		}
		return nil
	}
}
