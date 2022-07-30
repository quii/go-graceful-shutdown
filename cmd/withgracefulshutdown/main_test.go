package main

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/quii/go-graceful-shutdown/assert"
	"github.com/quii/go-graceful-shutdown/cmd"
)

const (
	url     = "http://localhost:8080"
	binName = "graceful"
)

func TestGracefulShutdown(t *testing.T) {
	deleteBinary, binPath, err := cmd.BuildBinary(binName)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(deleteBinary)

	sendInterrupt, err := cmd.RunServer(context.Background(), binPath)
	assert.NoError(t, err)

	// just check the server works before we shut things down
	_, err = http.Get(url)
	assert.NoError(t, err)

	// fire off a request, we know it is slow, and without graceful shutdown this would fail
	errCh := make(chan error, 1)
	go func() {
		_, err = http.Get(url)
		errCh <- err
	}()

	// give it a moment to fire the request and then send an interrupt
	time.Sleep(50 * time.Millisecond)
	assert.NoError(t, sendInterrupt())

	select {
	case err := <-errCh:
		assert.NoError(t, err)
	case <-time.After(3 * time.Second):
		t.Error("didnt get an error after 3s")
	}

	// after interrupt, the server should be shutdown, and no more requests will work
	_, err = http.Get(url)
	assert.Error(t, err)
}
