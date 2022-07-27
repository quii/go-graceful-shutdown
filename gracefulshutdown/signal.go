package gracefulshutdown

import (
	"os"
	"os/signal"
)

func NewInterruptSignalChannel() <-chan os.Signal {
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt)
	return osSignal
}
