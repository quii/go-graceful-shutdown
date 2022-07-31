package gracefulshutdown

import (
	"os"
	"os/signal"
	"syscall"
)

var (
	signalsToListenTo = []os.Signal{
		syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM,
	}
)

func newInterruptSignalChannel() <-chan os.Signal {
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, signalsToListenTo...)
	return osSignal
}
