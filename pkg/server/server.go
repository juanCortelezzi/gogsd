package server

import (
	"context"
	_ "embed"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/mattn/go-sqlite3"

	"github.com/juancortelezzi/gogsd/pkg/database"
	"github.com/juancortelezzi/gogsd/pkg/gsdlogger"
	"github.com/juancortelezzi/gogsd/pkg/routes"
)

func NewServerHandler(
	logger gsdlogger.Logger,
	queries *database.Queries,
	validate *validator.Validate,
) http.Handler {
	mux := http.NewServeMux()
	routes.AddRoutes(mux, logger, queries, validate)
	return mux
}

func Run(ctx context.Context, logger gsdlogger.Logger, lookupEnv func(string) (string, bool)) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	logger.DebugContext(ctx, "looking env variables")

	port, found := lookupEnv("PORT")
	if !found {
		return fmt.Errorf("PORT environment variable not found")
	}

	databaseUrl, found := lookupEnv("DATABASE_URL")
	if !found {
		return fmt.Errorf("DATABASE_URL environment variable not found")
	}

	logger.DebugContext(ctx, "initializing database conneciton")

	queries, err := database.Connect(ctx, logger, databaseUrl)
	if err != nil {
		logger.ErrorContext(ctx, "error connecting to database", "err", err)
		return nil
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	serverHandler := NewServerHandler(logger, queries, validate)

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
