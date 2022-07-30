package gracefulshutdown

import (
	"context"
	"net/http"
	"os"
	"time"
)

type (
	// HTTPServer is an abstraction of something that listens for connections and do HTTP things. 99% of the time, you'll pass in a net/http/Server.
	HTTPServer interface {
		ListenAndServe() error
		Shutdown(ctx context.Context) error
	}

	// Server wraps around a HTTPServer and will gracefully shutdown when it receives a shutdown signal.
	Server struct {
		shutdown <-chan os.Signal
		server   HTTPServer
		timeout  time.Duration
	}
)

// NewServer returns a Server, it allows you to send your own channel of signals you wish to shutdown gracefully on.
func NewServer(shutdown <-chan os.Signal, server HTTPServer, timeout time.Duration) *Server {
	return &Server{
		shutdown: shutdown,
		server:   server,
		timeout:  timeout,
	}
}

// NewDefaultServer wraps your HTTPServer with graceful shutdown against a "sensible" list of signals to listen to.
func NewDefaultServer(server HTTPServer, timeout time.Duration) *Server {
	return NewServer(NewInterruptSignalChannel(), server, timeout)
}

// Listen will call the ListenAndServe function of the HTTPServer you pass in. On a signal being sent to the shutdown signal provided in the constructor, it will call the server's Shutdown method to attempt to gracefully shutdown.
func (g *Server) Listen() error {
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
		ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
		defer cancel()

		// attempt to shutdown before ctx finishes (e.g a timeout)
		if err := g.server.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
			return err
		}
	}

	return nil
}
