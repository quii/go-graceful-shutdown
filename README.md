# my take on graceful shutdown of HTTP servers in Go

## The problem

- You're running a server, and deploying it many times per day
- You're also running tests
- Sometimes, you might be deploying a new version of the code while using the system, and if you're not doing this gracefully, you'll get some ugly 502 instead of your finished response

## The solution

Graceful shutdown! 

- Listen to interrupt signal
- Rather than killing the program straight away, instead call [http.Server.Shutdown](https://pkg.go.dev/net/http#Server.Shutdown) which will let requests, connections etc drain _before_ killing the server

There are a few examples of this out there, I thought I'd roll my own so i could understand it better and structure it in a non-confusing way, hopefully

## TODO

- I'd like to write some tests around it