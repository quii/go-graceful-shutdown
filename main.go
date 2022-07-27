package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	serverShutdownTimeout = 20 * time.Second
)

func main() {
	// create some kind of context with a timeout, so we don't just wait forever to shutdown
	applicationContext, _ := context.WithTimeout(context.Background(), serverShutdownTimeout)

	// notify of interrupt signals on a channel
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt)

	server := NewGracefulShutdownServer(
		osSignal,
		&http.Server{Addr: addr, Handler: http.HandlerFunc(aSlowHandler)},
	)

	if err := server.Listen(applicationContext); err != nil {
		// this will typically be if our responses aren't written before the ctx deadline, not much can be done
		log.Fatalf("uh oh, didnt shutdown gracefully, some responses may have been lost %v", err)
	}

	// hopefully, you'll always see this instead
	log.Println("shutdown gracefully! all responses were sent")
}
