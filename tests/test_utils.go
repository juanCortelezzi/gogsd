package tests

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/juancortelezzi/gogsd/pkg/gsdlogger"
)

func waitForReady(ctx context.Context, logger gsdlogger.Logger, timeout time.Duration, endpoint string) error {
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
				if time.Since(startTime) >= timeout {
					logger.ErrorContext(ctx, "error making request", "err", err)
					return fmt.Errorf("timeout reached while waiting for endpoint")
				}
				time.Sleep(time.Millisecond * 250)
				continue
			}
		}

		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			logger.InfoContext(ctx, "endpoint is ready!")
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if time.Since(startTime) >= timeout {
				return fmt.Errorf("timeout reached while waiting for endpoint")
			}
			time.Sleep(time.Millisecond * 250)
		}
	}
}

func testLookupEnv(key string) (string, bool) {
	if key == "PORT" {
		return "3000", true
	}
	return "", false
}
