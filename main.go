package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	applicationContext, _ := context.WithTimeout(context.Background(), serverShutdownTimeout)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	server := NewGracefulShutdownServer(c, http.HandlerFunc(aSlowHandler))

	if err := server.Listen(applicationContext); err != nil {
		log.Fatalf("uh oh, didnt shutdown gracefully, some responses may have been lost %v", err)
	}

	log.Println("shutdown gracefully! all responses were sent")
}
