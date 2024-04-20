package tests

import (
	"bufio"
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/juancortelezzi/gogsd/pkg/gsdlogger"
	"github.com/juancortelezzi/gogsd/pkg/server"
)

func TestHelloRoute(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	{
		logger := gsdlogger.NewLogger(os.Stdout, slog.LevelDebug)

		go server.Run(ctx, logger, testLookupEnv)

		err := waitForReady(ctx, logger, time.Second*3, getBaseUrl()+"/ping")
		if err != nil {
			t.Fatal(err)
		}
	}

	resp, err := http.Get(getBaseUrl() + "/hello/world")
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
