package gracefulshutdown

import (
	"context"
	"net/http"
	"os"
)

type (
	HTTPServer interface {
		ListenAndServe() error
		Shutdown(ctx context.Context) error
	}

	Server struct {
		shutdown <-chan os.Signal
		server   HTTPServer
	}
)

func NewServer(shutdown <-chan os.Signal, server HTTPServer) *Server {
	return &Server{
		shutdown: shutdown,
		server:   server,
	}
}

func NewDefaultServer(server HTTPServer) *Server {
	return NewServer(NewInterruptSignalChannel(), server)
}

func (g *Server) Listen(ctx context.Context) error {
	listenErr := make(chan error)

	// fly free, listen and serve
	go func() {
		if err := g.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			listenErr <- err
		}
	}()

	select {
	case err := <-listenErr:
		return err
	case <-g.shutdown:
		// attempt to shutdown before ctx finishes (e.g a timeout)
		if err := g.server.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
			return err
		}
	}

	return nil
}
