package main

import (
	"testing"
	"time"

	"github.com/quii/go-graceful-shutdown/acceptancetests"
	"github.com/quii/go-graceful-shutdown/assert"
)

const (
	port    = "8081"
	url     = "http://localhost:" + port
	binName = "without-graceful"
)

func TestNonGracefulShutdown(t *testing.T) {
	deleteBinary, binPath, err := acceptancetests.BuildBinary(binName)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(deleteBinary)

	sendInterrupt, err := acceptancetests.RunServer(binPath, port)
	assert.NoError(t, err)

	// just check the server works before we shut things down
	assert.NoError(t, acceptancetests.GetAndDiscardResponse(url))

	// fire off a request, it should fail because the server will be interrupted
	time.AfterFunc(50*time.Millisecond, func() {
		assert.NoError(t, sendInterrupt())
	})
	errCh := make(chan error, 1)
	go func() {
		errCh <- acceptancetests.GetAndDiscardResponse(url)
	}()

	select {
	case err := <-errCh:
		assert.Error(t, err)
	case <-time.After(3 * time.Second):
		t.Error("didnt get an error after 3s")
	}

	// after interrupt, the server should be shutdown, and no more requests will work
	assert.Error(t, acceptancetests.GetAndDiscardResponse(url))
}
