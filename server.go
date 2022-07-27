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

type GracefulShutdownServer struct {
	shutdown <-chan os.Signal
	server   http.Server
}

func NewGracefulShutdownServer(shutdown <-chan os.Signal, handler http.Handler) *GracefulShutdownServer {
	return &GracefulShutdownServer{
		shutdown: shutdown,
		server:   http.Server{Addr: addr, Handler: handler},
	}
}

func (g *GracefulShutdownServer) Listen(ctx context.Context) error {
	// fly free, listen and serve
	go func() {
		if err := g.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
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
