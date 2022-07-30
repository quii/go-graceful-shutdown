package cmd

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

const (
	baseBinName = "temp-testbinary"
)

func BuildBinary(name string) (cleanup func(), cmdPath string, err error) {
	binName := name + "-" + baseBinName

	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	build := exec.Command("go", "build", "-o", binName)

	if err := build.Run(); err != nil {
		return nil, "", fmt.Errorf("cannot build tool %s: %s", binName, err)
	}

	dir, err := os.Getwd()
	if err != nil {
		return nil, "", err
	}

	cmdPath = filepath.Join(dir, binName)

	cleanup = func() {
		os.Remove(binName)
	}

	return
}

func RunServer(ctx context.Context, path string) (sendInterrupt func() error, err error) {
	cmd := exec.CommandContext(ctx, path)
	cmd.Stderr = NewLogWriter()

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("cannot run temp converter: %s", err)
	}
	waitForServerListening()

	sendInterrupt = func() error {
		return cmd.Process.Signal(os.Interrupt)
	}

	return sendInterrupt, nil
}

func waitForServerListening() {
	for i := 0; i < 20; i++ {
		conn, _ := net.Dial("tcp", net.JoinHostPort("localhost", "8080"))
		if conn != nil {
			conn.Close()
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}

type LogWriter struct {
	logger *log.Logger
}

func NewLogWriter() *LogWriter {
	lw := &LogWriter{}
	lw.logger = log.Default()
	return lw
}

func (lw LogWriter) Write(p []byte) (n int, err error) {
	lw.logger.Println(string(p))
	return len(p), nil
}
