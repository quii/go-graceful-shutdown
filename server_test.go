package gracefulshutdown_test

import (
	"context"
	"errors"
	"github.com/quii/graceful-shutdown"
	"github.com/quii/graceful-shutdown/assert"
	"os"
	"testing"
	"time"
)

func TestGracefulShutdownServer_Listen(t *testing.T) {
	t.Run("happy path, listen, wait for interrupt, shutdown gracefully", func(t *testing.T) {
		var (
			interrupt = make(chan os.Signal)
			spyServer = gracefulshutdown.NewSpyServer()
			server    = gracefulshutdown.NewServer(interrupt, spyServer)
		)

		spyServer.ListenAndServeFunc = func() error {
			return nil
		}
		spyServer.ShutdownFunc = func() error {
			return nil
		}

		go server.Listen(context.Background())

		// verify we call listen on the delegate server
		spyServer.AssertListened(t)

		// verify we call shutdown on the delegate server when an interrupt is made
		interrupt <- os.Interrupt
		spyServer.AssertShutdown(t)
	})

	t.Run("when listen fails, return error", func(t *testing.T) {
		var (
			interrupt = make(chan os.Signal)
			spyServer = gracefulshutdown.NewSpyServer()
			server    = gracefulshutdown.NewServer(interrupt, spyServer)
			err       = errors.New("oh no")
		)

		spyServer.ListenAndServeFunc = func() error {
			return err
		}

		gotErr := server.Listen(context.Background())

		assert.Equal(t, gotErr.Error(), err.Error())
	})

	t.Run("shutdown error gets propagated", func(t *testing.T) {
		var (
			interrupt = make(chan os.Signal)
			errChan   = make(chan error)
			spyServer = gracefulshutdown.NewSpyServer()
			server    = gracefulshutdown.NewServer(interrupt, spyServer)
			err       = errors.New("oh no")
		)

		spyServer.ListenAndServeFunc = func() error {
			return nil
		}
		spyServer.ShutdownFunc = func() error {
			return err
		}

		go func() {
			errChan <- server.Listen(context.Background())
		}()

		interrupt <- os.Interrupt

		select {
		case gotErr := <-errChan:
			assert.Equal(t, gotErr.Error(), err.Error())
		case <-time.After(500 * time.Millisecond):
			t.Error("timed out waiting for shutdown error to be propagated")
		}
	})

}
