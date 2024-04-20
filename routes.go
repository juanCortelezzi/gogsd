package main

import (
	"fmt"
	"net/http"
	"time"
)

func debugMiddleware(logger Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			routeLogger := logger.With(
				"method", r.Method,
				"route", r.URL.Path,
			)
			routeLogger.InfoContext(r.Context(), "hitStart")
			now := time.Now()
			h.ServeHTTP(w, r)
			routeLogger.InfoContext(r.Context(), "hitEnd", "timeTaken", time.Since(now))
		})
	}
}

func handlePing() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Pong")
	})
}

func handleHello(logger Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")

		if name == "" {
			logger.ErrorContext(r.Context(), "error getting name from url", "name", name)
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, "Hello, %s!\n", name)
	})
}

func AddRoutes(mux *http.ServeMux, logger Logger) {
	debugMiddle := debugMiddleware(logger)

	mux.Handle("GET /ping", handlePing())
	mux.Handle("GET /hello/{name}", debugMiddle(handleHello(logger)))
}
