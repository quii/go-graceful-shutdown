# Graceful shutdown decorator
[![Go Reference](https://pkg.go.dev/badge/github.com/quii/go-graceful-shutdown.svg)](https://pkg.go.dev/github.com/quii/go-graceful-shutdown)
![Pipeline](https://github.com/quii/go-graceful-shutdown/actions/workflows/pipeline.yaml/badge.svg)
![Lint](https://github.com/quii/go-graceful-shutdown/actions/workflows/golang-ci-lint.yaml/badge.svg)
```go
func main() {
	httpServer := &http.Server{Addr: ":8080", Handler: http.HandlerFunc(aSlowHandler)}

	server := gracefulshutdown.NewDefaultServer(httpServer, serverShutdownTimeout)

	if err := server.Listen(); err != nil {
		// this will typically happen if our responses aren't written before the ctx deadline, not much can be done
		log.Fatalf("uh oh, didnt shutdown gracefully, some responses may have been lost %v", err)
	}

	// hopefully, you'll always see this instead
	log.Println("shutdown gracefully! all responses were sent")
}
```

## The problem

- You're running a HTTP server, and deploying it many times per day
- Sometimes, you might be deploying a new version of the code while it is trying to handle a request, and if you're not handling this gracefully you'll either:
  - Not get a response
  - Or the reverse-proxy in front of your service will complain about your service and return a 502

## The solution

Graceful shutdown! 

- Listen to interrupt signals
- Rather than killing the program straight away, instead call [http.Server.Shutdown](https://pkg.go.dev/net/http#Server.Shutdown) which will let requests, connections e.t.c drain _before_ killing the server
- This should mean in most cases, the server will finish the currently running requests before stopping

There are a few examples of this out there, I thought I'd roll my own so I could understand it better, and structure it in a non-confusing way, hopefully.

Almost everything boils down to a decorator pattern in the end. You provide my library a `*http.Server` and it'll return you back a `*gracefulshutdown.Server`. Just call `Listen` instead of your normal `ListenAndServe`, and it'll gracefully shutdown on [an os signal](https://github.com/quii/go-graceful-shutdown/blob/main/signal.go#L11).

## Example usage

See [cmd/main.go](https://github.com/quii/go-graceful-shutdown/blob/main/cmd/main.go)