package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func NewServerHandler(logger Logger) http.Handler {
	mux := http.NewServeMux()
	AddRoutes(mux, logger)
	return mux
}

func Run(ctx context.Context, logger Logger, lookupEnv func(string) (string, bool)) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	serverHandler := NewServerHandler(logger)

	port, found := lookupEnv("PORT")
	if !found {
		return fmt.Errorf("PORT environment variable not found")
	}

	httpServer := &http.Server{
		Addr:    net.JoinHostPort("127.0.0.1", port),
		Handler: serverHandler,
	}

	go func() {
		logger.InfoContext(ctx, "listening on", "addr", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.ErrorContext(ctx, "error listening and serving", "err", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.ErrorContext(ctx, "error shutting down http server", "err", err)
		}
	}()

	wg.Wait()

	return nil
}
