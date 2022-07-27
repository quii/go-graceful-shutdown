package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	serverShutdownTimeout = 20 * time.Second
	addr                  = ":8080"
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
	go func() {
		if err := g.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	osCall := <-g.shutdown
	log.Printf("system call: %+v", osCall)

	err := g.server.Shutdown(ctx)
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}
