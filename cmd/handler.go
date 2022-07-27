package main

import (
	"fmt"
	"net/http"
	"time"
)

func aSlowHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second)
	fmt.Fprint(w, "Hello, world")
}
