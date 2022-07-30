package main

import (
	"context"
	"log"
	"net/http"
	"time"

	gracefulshutdown "github.com/quii/go-graceful-shutdown"
)

const (
	serverShutdownTimeout = 20 * time.Second
	addr                  = ":8080"
)

func main() {
	// create a context with a timeout, so we don't just wait forever to shut down
	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)

	httpServer := &http.Server{Addr: addr, Handler: http.HandlerFunc(aSlowHandler)}

	server := gracefulshutdown.NewDefaultServer(httpServer)

	if err := server.Listen(ctx); err != nil {
		// this will typically happen if our responses aren't written before the ctx deadline, not much can be done
		log.Fatalf("uh oh, didnt shutdown gracefully, some responses may have been lost %v", err)
	}

	cancel()

	// hopefully, you'll always see this instead
	log.Println("shutdown gracefully! all responses were sent")
}
