package main

import (
	"testing"
	"time"

	"github.com/quii/go-graceful-shutdown/acceptancetests"
	"github.com/quii/go-graceful-shutdown/assert"
)

const (
	port    = "8080"
	url     = "http://localhost:" + port
	binName = "graceful"
)

func TestGracefulShutdown(t *testing.T) {
	deleteBinary, binPath, err := acceptancetests.BuildBinary(binName)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(deleteBinary)

	sendInterrupt, err := acceptancetests.RunServer(binPath, port)
	assert.NoError(t, err)

	// just check the server works before we shut things down
	assert.NoError(t, acceptancetests.GetAndDiscardResponse(url))

	// fire off a request, we know it is slow, and without graceful shutdown this would fail
	time.AfterFunc(50*time.Millisecond, func() {
		assert.NoError(t, sendInterrupt())
	})
	errCh := make(chan error, 1)
	go func() {
		errCh <- acceptancetests.GetAndDiscardResponse(url)
	}()

	select {
	case err := <-errCh:
		assert.NoError(t, err)
	case <-time.After(3 * time.Second):
		t.Error("didnt get an error after 3s")
	}

	// after interrupt, the server should be shutdown, and no more requests will work
	assert.Error(t, acceptancetests.GetAndDiscardResponse(url))
}
