package main

import (
	"log"
	"net/http"

	"github.com/quii/go-graceful-shutdown/cmd"
)

func main() {
	server := &http.Server{Addr: ":8081", Handler: http.HandlerFunc(cmd.ASlowHandler)}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
