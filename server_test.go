package gracefulshutdown_test

import (
	"context"
	"errors"
	"github.com/quii/graceful-shutdown"
	"github.com/quii/graceful-shutdown/assert"
	"os"
	"testing"
)

func TestGracefulShutdownServer_Listen(t *testing.T) {
	t.Run("happy path, listen, wait for interrupt, shutdown gracefully", func(t *testing.T) {

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
	})

	t.Run("when listen fails, return error", func(t *testing.T) {
		interrupt := make(chan os.Signal)
		spyServer := gracefulshutdown.NewSpyServer()
		err := errors.New("oh no")
		spyServer.ListenAndServeFunc = func() error {
			return err
		}
		server := gracefulshutdown.NewServer(interrupt, spyServer)
		gotErr := server.Listen(context.Background())

		assert.Equal(t, gotErr.Error(), err.Error())
	})
	
	t.Skip("shutdown error gets propagated")

}
