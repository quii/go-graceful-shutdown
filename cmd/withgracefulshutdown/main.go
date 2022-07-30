package main

import (
	"log"
	"net/http"
	"time"

	gracefulshutdown "github.com/quii/go-graceful-shutdown"
	"github.com/quii/go-graceful-shutdown/cmd"
)

const (
	serverShutdownTimeout = 10 * time.Second
)

func main() {
	httpServer := &http.Server{Addr: ":8080", Handler: http.HandlerFunc(cmd.ASlowHandler)}

	server := gracefulshutdown.NewDefaultServer(httpServer, serverShutdownTimeout)

	if err := server.Listen(); err != nil {
		// this will typically happen if our responses aren't written before the ctx deadline, not much can be done
		log.Fatalf("uh oh, didnt shutdown gracefully, some responses may have been lost %v", err)
	}

	// hopefully, you'll always see this instead
	log.Println("shutdown gracefully! all responses were sent")
}
