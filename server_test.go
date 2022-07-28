package gracefulshutdown_test

import (
	"context"
	"github.com/quii/graceful-shutdown"
	"os"
	"testing"
)

func TestGracefulShutdownServer_Listen(t *testing.T) {
	interrupt := make(chan os.Signal)
	spyServer := gracefulshutdown.NewSpyServer()
	spyServer.ListenAndServeFunc = func() error {
		return nil
	}
	spyServer.ShutdownFunc = func() error {
		return nil
	}

	server := gracefulshutdown.NewServer(interrupt, spyServer)
	go server.Listen(context.Background())

	// verify we call listen on the delegate server
	spyServer.AssertListened(t)

	// verify we call shutdown on the delegate server when an interrupt is made
	interrupt <- os.Interrupt
	spyServer.AssertShutdown(t)
}
