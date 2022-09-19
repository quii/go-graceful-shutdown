# Graceful shutdown decorator
[![Go Reference](https://pkg.go.dev/badge/github.com/quii/go-graceful-shutdown.svg)](https://pkg.go.dev/github.com/quii/go-graceful-shutdown)
![Pipeline](https://github.com/quii/go-graceful-shutdown/actions/workflows/pipeline.yaml/badge.svg)
![Lint](https://github.com/quii/go-graceful-shutdown/actions/workflows/golang-ci-lint.yaml/badge.svg)

A wrapper for your Go HTTP server so that it will finish responding to in-flight requests on interrupt signals before shutting down.

```go
func main() {
  var (
    ctx        = context.Background()
    httpServer = &http.Server{Addr: ":8080", Handler: http.HandlerFunc(acceptancetests.SlowHandler)}
    server     = gracefulshutdown.NewServer(httpServer)
  )

  if err := server.ListenAndServe(ctx); err != nil {
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

There are a few examples of this out there, I thought I'd roll my own, so I could understand it better, and structure it in a non-confusing way, hopefully.

Almost everything boils down to a decorator pattern in the end. You provide my library a `*http.Server` and it'll return you back a `*gracefulshutdown.Server`. Just call `ListenAndServe`, and it'll gracefully shutdown on [an os signal](https://github.com/quii/go-graceful-shutdown/blob/main/signal.go#L11).

## Example usage and testing

See [acceptancetests/withgracefulshutdown/main.go](https://github.com/quii/go-graceful-shutdown/blob/main/acceptancetests/withgracefulshutdown/main.go) for an example

There are two binaries in this project with accompanying acceptance tests to verify the functionality that live inside `/acceptancetests`.

Both tests build the binaries, run them, fire a `HTTP GET` and then send an interrupt signal to tell the server to stop.

The two binaries allow us to test both scenarios

1. A "slow" HTTP server with no graceful shutdown. For this we assert that we do get an error, because the server should shutdown immediately and any in-flight requests will fail.
2. Another slow HTTP server _with_ graceful shutdown. Same test again, but this time we assert we don't get an error as we expect to get a response before the server is terminated.