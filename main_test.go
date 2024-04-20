package main

import (
	"bufio"
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"
)

const (
	baseUrl = "http://127.0.0.1:3000"
)

func TestHelloRoute(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	{
		logger := NewLogger(os.Stdout, slog.LevelDebug)

		go Run(ctx, logger, testLookupEnv)

		err := WaitForReady(ctx, logger, time.Second*3, baseUrl+"/ping")
		if err != nil {
			t.Fatal(err)
		}
	}

	resp, err := http.Get(baseUrl + "/hello/world")
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	data, err := bufio.NewReaderSize(resp.Body, 1024).ReadBytes('\n')
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(data, []byte("Hello, world!\n")) {
		t.Fatalf("expected Hello, world! got %s\n", data)
	}
}
