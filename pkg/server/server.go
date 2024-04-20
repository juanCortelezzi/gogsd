package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/juancortelezzi/gogsd/pkg/gsdlogger"
	"github.com/juancortelezzi/gogsd/pkg/routes"
)

func NewServerHandler(l gsdlogger.Logger) http.Handler {
	mux := http.NewServeMux()
	routes.AddRoutes(mux, l)
	return mux
}

func Run(ctx context.Context, l gsdlogger.Logger, lookupEnv func(string) (string, bool)) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	serverHandler := NewServerHandler(l)

	port, found := lookupEnv("PORT")
	if !found {
		return fmt.Errorf("PORT environment variable not found")
	}

	httpServer := &http.Server{
		Addr:    net.JoinHostPort("127.0.0.1", port),
		Handler: serverHandler,
	}

	go func() {
		l.InfoContext(ctx, "listening on", "addr", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.ErrorContext(ctx, "error listening and serving", "err", err)
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
			l.ErrorContext(ctx, "error shutting down http server", "err", err)
		}
	}()

	wg.Wait()

	return nil
}
