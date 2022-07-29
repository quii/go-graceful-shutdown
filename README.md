# my take on graceful shutdown of HTTP servers in Go

## The problem

- You're running a server, and deploying it many times per day
- You're also running tests
- Sometimes, you might be deploying a new version of the code while using the system, and if you're not doing this gracefully you'll either:
  - Just not get a response
  - Or the reverse-proxy in front of your service will return a 502

## The solution

Graceful shutdown! 

- Listen to interrupt signal
- Rather than killing the program straight away, instead call [http.Server.Shutdown](https://pkg.go.dev/net/http#Server.Shutdown) which will let requests, connections etc drain _before_ killing the server
- This should mean in most cases, the server will finish the currently running requests before stopping

There are a few examples of this out there, I thought I'd roll my own so I could understand it better, and structure it in a non-confusing way, hopefully.

Almost everything boils down to a decorator pattern in the end. You provide my library a `*http.Server` and it'll return you back a `*gracefulshutdown.Server`. Just call `Listen` instead of your normal `ListenAndServe`, and it'll gracefully shutdown on [an os signal](https://github.com/quii/go-graceful-shutdown/blob/main/signal.go#L11).

## Example usage

See [cmd/main.go](https://github.com/quii/go-graceful-shutdown/blob/main/cmd/main.go)