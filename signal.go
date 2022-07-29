package gracefulshutdown

import (
	"os"
	"os/signal"
	"syscall"
)

var (
	// SignalsToListenTo is a sensible default list of signals from the OS to listen to.
	SignalsToListenTo = []os.Signal{
		syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM,
	}
)

// NewInterruptSignalChannel returns a channel which will be notified on any of the SignalsToListenTo.
func NewInterruptSignalChannel() <-chan os.Signal {
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, SignalsToListenTo...)
	return osSignal
}
