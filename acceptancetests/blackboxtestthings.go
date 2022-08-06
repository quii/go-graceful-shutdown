package acceptancetests

import (
	"fmt"
	"log"
	"math/rand"
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

func LaunchTestProgram(port string) (cleanup func(), sendInterrupt func() error, err error) {
	binName := randomString(10) + "-" + baseBinName

	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	build := exec.Command("go", "build", "-o", binName)

	if err := build.Run(); err != nil {
		return nil, nil, fmt.Errorf("cannot build tool %s: %s", binName, err)
	}

	dir, err := os.Getwd()
	if err != nil {
		return nil, nil, err
	}

	cmdPath := filepath.Join(dir, binName)

	sendInterrupt, kill, err := runServer(cmdPath, port)

	cleanup = func() {
		if kill != nil {
			kill()
		}
		os.Remove(binName)
	}

	if err != nil {
		cleanup() // even though it's not listening correctly, the program could still be running
		return nil, nil, err
	}

	return cleanup, sendInterrupt, nil
}

func runServer(path string, port string) (sendInterrupt func() error, kill func(), err error) {
	cmd := exec.Command(path)
	cmd.Stderr = NewLogWriter()

	if err := cmd.Start(); err != nil {
		return nil, nil, fmt.Errorf("cannot run temp converter: %s", err)
	}

	kill = func() {
		_ = cmd.Process.Kill()
	}

	sendInterrupt = func() error {
		return cmd.Process.Signal(os.Interrupt)
	}

	err = waitForServerListening(port)

	return
}

func waitForServerListening(port string) error {
	for i := 0; i < 30; i++ {
		conn, _ := net.Dial("tcp", net.JoinHostPort("localhost", port))
		if conn != nil {
			conn.Close()
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("nothing seems to be listening on localhost:%s", port)
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

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
