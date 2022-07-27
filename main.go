package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type GracefulShutdownServer struct {
	shutdown <-chan os.Signal
	server   http.Server
}

func NewGracefulShutdownServer(shutdown <-chan os.Signal) *GracefulShutdownServer {
	server := http.Server{Addr: ":8080"}
	server.Handler = http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		time.Sleep(5 * time.Second)
		fmt.Fprint(w, "Hello, world")
	})

	return &GracefulShutdownServer{shutdown: shutdown, server: server}
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

func main() {
	// create an application context we can use
	// the timeout on the context gives Shutdown its deadline to finish responding with requests
	applicationContext, _ := context.WithTimeout(context.Background(), 20*time.Second)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	server := NewGracefulShutdownServer(c)

	if err := server.Listen(applicationContext); err != nil {
		log.Fatalf("uh oh, didnt shutdown gracefully, some responses may have been lost %v", err)
	}

	log.Println("shutdown gracefully! all responses were sent")
}
