package main

import (
	"context"
	"net/http"
	"os"
)

const (
	addr = ":8080"
)

type (
	Server interface {
		ListenAndServe() error
		Shutdown(ctx context.Context) error
	}

	GracefulShutdownServer struct {
		shutdown <-chan os.Signal
		server   Server
	}
)

func NewGracefulShutdownServer(shutdown <-chan os.Signal, server Server) *GracefulShutdownServer {
	return &GracefulShutdownServer{
		shutdown: shutdown,
		server:   server,
	}
}

func (g *GracefulShutdownServer) Listen(ctx context.Context) error {
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
