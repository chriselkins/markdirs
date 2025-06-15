package main

import (
	"bytes"
	"flag"
	"io"
	"os"
	"strings"
	"testing"
)

func captureStderr(f func()) string {
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	done := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		done <- buf.String()
	}()

	f()
	w.Close()
	os.Stderr = old
	return <-done
}

func TestRun_BadArgs_ShowsUsage(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	os.Args = []string{"markdirs", "-q", "/tmp", ".file"}

	var code int
	output := captureStderr(func() {
		code = run()
	})

	if code != 1 {
		t.Errorf("Expected exit code 1, got %d", code)
	}

	if !strings.Contains(output, "Usage:") {
		t.Errorf("Expected usage message, got %q", output)
	}
}
