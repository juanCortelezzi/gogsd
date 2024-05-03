package main

import (
	"context"
	"log/slog"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/juancortelezzi/gogsd/pkg/gsdlogger"
	"github.com/juancortelezzi/gogsd/pkg/server"
)

func main() {
	ctx := context.Background()
	logger := gsdlogger.NewLogger(os.Stdout, slog.LevelInfo)

	if err := server.Run(ctx, logger, os.LookupEnv); err != nil {
		logger.ErrorContext(ctx, "error in top level", "err", err)
		os.Exit(1)
	}
}
