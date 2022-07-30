package gracefulshutdown_test

import (
	"context"
	"testing"

	"github.com/quii/go-graceful-shutdown/assert"
)

type SpyServer struct {
	ListenAndServeFunc func() error
	listened           chan struct{}

	ShutdownFunc func() error
	shutdown     chan struct{}
}

func NewSpyServer() *SpyServer {
	return &SpyServer{
		listened: make(chan struct{}, 1),
		shutdown: make(chan struct{}, 1),
	}
}

func (s *SpyServer) ListenAndServe() error {
	s.listened <- struct{}{}
	return s.ListenAndServeFunc()
}

func (s *SpyServer) AssertListened(t *testing.T) {
	t.Helper()
	assert.SignalSent(t, s.listened, "listen")
}

func (s *SpyServer) Shutdown(ctx context.Context) error {
	s.shutdown <- struct{}{}
	return s.ShutdownFunc()
}

func (s *SpyServer) AssertShutdown(t *testing.T) {
	t.Helper()
	assert.SignalSent(t, s.shutdown, "shutdown")
}
