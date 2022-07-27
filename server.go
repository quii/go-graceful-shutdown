package main

import (
	"context"
	"log"
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
	// fly free, listen and serve
	go func() {
		if err := g.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
			//todo: i think this might be buggy, this would kill the go routine but not the server, which is what we probably want here
		}
	}()

	// wait for the shutdown signal
	<-g.shutdown

	// attempt to shutdown before ctx finishes (e.g a timeout)
	if err := g.server.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
