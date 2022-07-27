package main

import (
	"context"
	"github.com/quii/graceful-shutdown/gracefulshutdown"
	"log"
	"net/http"
	"time"
)

const (
	serverShutdownTimeout = 20 * time.Second
	addr                  = ":8080"
)

func main() {
	// create some kind of context with a timeout, so we don't just wait forever to shutdown
	applicationContext, _ := context.WithTimeout(context.Background(), serverShutdownTimeout)

	server := gracefulshutdown.NewDefaultServer(
		&http.Server{Addr: addr, Handler: http.HandlerFunc(aSlowHandler)},
	)

	if err := server.Listen(applicationContext); err != nil {
		// this will typically be if our responses aren't written before the ctx deadline, not much can be done
		log.Fatalf("uh oh, didnt shutdown gracefully, some responses may have been lost %v", err)
	}

	// hopefully, you'll always see this instead
	log.Println("shutdown gracefully! all responses were sent")
}
