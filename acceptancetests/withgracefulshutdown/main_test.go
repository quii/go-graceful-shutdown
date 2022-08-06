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

	sendInterrupt, kill, err := acceptancetests.RunServer(binPath, port)
	t.Cleanup(kill)
	assert.NoError(t, err)

	// just check the server works before we shut things down
	assert.CanGet(t, url)

	// fire off a request, we know it is slow, and without graceful shutdown this would fail
	time.AfterFunc(50*time.Millisecond, func() {
		assert.NoError(t, sendInterrupt())
	})
	assert.CanGet(t, url)

	// after interrupt, the server should be shutdown, and no more requests will work
	assert.CantGet(t, url)
}
