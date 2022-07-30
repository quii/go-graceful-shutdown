package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/quii/go-graceful-shutdown/assert"
	"github.com/quii/go-graceful-shutdown/cmd"
)

const (
	port    = "8081"
	url     = "http://localhost:" + port
	binName = "without-graceful"
)

func TestNonGracefulShutdown(t *testing.T) {
	deleteBinary, binPath, err := cmd.BuildBinary(binName)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(deleteBinary)

	sendInterrupt, err := cmd.RunServer(binPath, port)
	assert.NoError(t, err)

	// just check the server works before we shut things down
	_, err = http.Get(url)
	assert.NoError(t, err)

	// fire off a request, it should fail because the server will be interrupted
	time.AfterFunc(50*time.Millisecond, func() {
		assert.NoError(t, sendInterrupt())
	})
	errCh := make(chan error, 1)
	go func() {
		_, err = http.Get(url)
		errCh <- err
	}()

	select {
	case err := <-errCh:
		assert.Error(t, err)
	case <-time.After(3 * time.Second):
		t.Error("didnt get an error after 3s")
	}

	// after interrupt, the server should be shutdown, and no more requests will work
	_, err = http.Get(url)
	assert.Error(t, err)
}
