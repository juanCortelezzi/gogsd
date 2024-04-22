package tests

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/juancortelezzi/gogsd/pkg/gsdlogger"
)

const waitForReadyTimeout = time.Second * 3

func testLookupEnv(key string) (string, bool) {
	switch key {
	case "PORT":
		return "3000", true
	case "DATABASE_URL":
		return ":memory:", true
	default:
		return "", false
	}
}

func getBaseUrl() string {
	port, found := testLookupEnv("PORT")
	if !found {
		panic("PORT environment variable not found")
	}

	return "http://" + net.JoinHostPort("127.0.0.1", port)
}

func waitForReady(ctx context.Context, logger gsdlogger.Logger, endpoint string) error {
	client := http.Client{}
	startTime := time.Now()

	for {

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		if err != nil {
			return fmt.Errorf("failed to create a request: %w", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				if time.Since(startTime) >= waitForReadyTimeout {
					logger.ErrorContext(ctx, "error making request", "err", err)
					return fmt.Errorf("timeout reached while waiting for endpoint")
				}
				time.Sleep(time.Millisecond * 250)
				continue
			}
		}

		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if time.Since(startTime) >= waitForReadyTimeout {
				return fmt.Errorf("timeout reached while waiting for endpoint")
			}
			time.Sleep(time.Millisecond * 250)
		}
	}
}
