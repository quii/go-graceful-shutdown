package gracefulshutdown

import (
	"context"
	"net/http"
	"os"
	"time"
)

const (
	k8sDefaultTerminationGracePeriod = 30 * time.Second
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

	// ServerOption provides ways of configuring Server.
	ServerOption func(server *Server)
)

// WithShutdownSignal WithShutdownSignals allows you to listen to whatever signals you like, rather than the default ones defined in signal.go.
func WithShutdownSignal(shutdownSignal <-chan os.Signal) ServerOption {
	return func(server *Server) {
		server.shutdown = shutdownSignal
	}
}

// WithTimeout lets you set your own timeout for waiting for graceful shutdown. By default this is set to 30 seconds (k8s' default TerminationGracePeriod).
func WithTimeout(timeout time.Duration) ServerOption {
	return func(server *Server) {
		server.timeout = timeout
	}
}

// NewServer returns a Server that can gracefully shutdown on shutdown signals.
func NewServer(server HTTPServer, options ...ServerOption) *Server {
	s := &Server{
		delegate: server,
		timeout:  k8sDefaultTerminationGracePeriod,
		shutdown: newInterruptSignalChannel(),
	}

	for _, option := range options {
		option(s)
	}

	return s
}

// ListenAndServe will call the ListenAndServe function of the delegate HTTPServer you passed in at construction. On a signal being sent to the shutdown signal provided in the constructor, it will call the server's Shutdown method to attempt to gracefully shutdown.
func (s *Server) ListenAndServe(ctx context.Context) error {
	select {
	case err := <-s.delegateListenAndServe():
		return err
	case <-ctx.Done():
		return s.shutdownDelegate(ctx)
	case <-s.shutdown:
		return s.shutdownDelegate(ctx)
	}
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

func (s *Server) shutdownDelegate(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if err := s.delegate.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		return err
	}
	return ctx.Err()
}
