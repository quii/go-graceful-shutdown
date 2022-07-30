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
		delegate HTTPServer
		timeout  time.Duration
	}
)

// NewServer returns a Server, it allows you to send your own channel of signals you wish to shutdown gracefully on.
func NewServer(shutdown <-chan os.Signal, server HTTPServer, timeout time.Duration) *Server {
	return &Server{
		shutdown: shutdown,
		delegate: server,
		timeout:  timeout,
	}
}

// NewDefaultServer wraps your HTTPServer with graceful shutdown against a "sensible" list of signals to listen to.
func NewDefaultServer(server HTTPServer, timeout time.Duration) *Server {
	return NewServer(NewInterruptSignalChannel(), server, timeout)
}

// ListenAndServe will call the ListenAndServe function of the delegate HTTPServer you passed in at construction. On a signal being sent to the shutdown signal provided in the constructor, it will call the server's Shutdown method to attempt to gracefully shutdown.
func (s *Server) ListenAndServe() error {
	select {
	case err := <-s.delegateListenAndServe():
		return err
	case <-s.shutdown:
		if err := s.shutdownDelegate(); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) delegateListenAndServe() chan error {
	listenErr := make(chan error)

	go func() {
		if err := s.delegate.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			listenErr <- err
		}
	}()

	return listenErr
}

func (s *Server) shutdownDelegate() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	if err := s.delegate.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}
