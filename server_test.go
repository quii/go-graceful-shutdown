package gracefulshutdown_test

import (
	"errors"
	"os"
	"testing"
	"time"

	gracefulshutdown "github.com/quii/go-graceful-shutdown"
	"github.com/quii/go-graceful-shutdown/assert"
)

const (
	timeout = 5 * time.Second
)

func TestGracefulShutdownServer_Listen(t *testing.T) {
	t.Run("happy path, listen, wait for interrupt, shutdown gracefully", func(t *testing.T) {
		var (
			interrupt = make(chan os.Signal)
			spyServer = NewSpyServer()
			server    = gracefulshutdown.NewServer(interrupt, spyServer, timeout)
		)

		spyServer.ListenAndServeFunc = func() error {
			return nil
		}
		spyServer.ShutdownFunc = func() error {
			return nil
		}

		go func() {
			if err := server.Listen(); err != nil {
				t.Error(err)
			}
		}()

		// verify we call listen on the delegate server
		spyServer.AssertListened(t)

		// verify we call shutdown on the delegate server when an interrupt is made
		interrupt <- os.Interrupt
		spyServer.AssertShutdown(t)
	})

	t.Run("when listen fails, return error", func(t *testing.T) {
		var (
			interrupt = make(chan os.Signal)
			spyServer = NewSpyServer()
			server    = gracefulshutdown.NewServer(interrupt, spyServer, timeout)
			err       = errors.New("oh no")
		)

		spyServer.ListenAndServeFunc = func() error {
			return err
		}

		gotErr := server.Listen()

		assert.Equal(t, gotErr.Error(), err.Error())
	})

	t.Run("shutdown error gets propagated", func(t *testing.T) {
		var (
			interrupt = make(chan os.Signal)
			errChan   = make(chan error)
			spyServer = NewSpyServer()
			server    = gracefulshutdown.NewServer(interrupt, spyServer, timeout)
			err       = errors.New("oh no")
		)

		spyServer.ListenAndServeFunc = func() error {
			return nil
		}
		spyServer.ShutdownFunc = func() error {
			return err
		}

		go func() {
			errChan <- server.Listen()
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
