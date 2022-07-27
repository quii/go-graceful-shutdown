package gracefulshutdown_test

import (
	"context"
	"github.com/quii/graceful-shutdown"
	"github.com/quii/graceful-shutdown/assert"
	"os"
	"testing"
	"time"
)

type SpyServer struct {
	ListenCalls        int
	ListenAndServeFunc func() error

	ShutdownCalls int
	ShutdownFunc  func() error
}

func (s *SpyServer) ListenAndServe() error {
	s.ListenCalls++
	return s.ListenAndServeFunc()
}

func (s *SpyServer) Shutdown(ctx context.Context) error {
	s.ShutdownCalls++
	return s.ShutdownFunc()
}

//todo: this is shite but better than nothing
func TestGracefulShutdownServer_Listen(t *testing.T) {
	interrupt := make(chan os.Signal)
	serverSpy := &SpyServer{
		ListenAndServeFunc: func() error {
			return nil
		},
		ShutdownFunc: func() error {
			return nil
		},
	}

	server := gracefulshutdown.NewServer(interrupt, serverSpy)
	go server.Listen(context.Background())

	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, serverSpy.ListenCalls, 1)

	interrupt <- os.Interrupt
	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, serverSpy.ShutdownCalls, 1)
}
