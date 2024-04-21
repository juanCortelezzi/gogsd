package routes

import (
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/juancortelezzi/gogsd/pkg/database"
	"github.com/juancortelezzi/gogsd/pkg/gsdlogger"
	"github.com/juancortelezzi/gogsd/pkg/handlers"
)

func AddRoutes(
	mux *http.ServeMux,
	logger gsdlogger.Logger,
	queries *database.Queries,
	validate *validator.Validate,
) {
	debugMiddle := debugMiddleware(logger)

	mux.Handle("GET /ping", handlers.HandlePing())
	mux.Handle("GET /hello/{name}", debugMiddle(handlers.HandleHello(logger)))

	mux.Handle("GET /todos", debugMiddle(handlers.HandleListTodos(logger, queries)))
	mux.Handle("POST /todos", debugMiddle(handlers.HandleCreateTodo(logger, queries, validate)))
	mux.Handle("PUT /todos/{id}", debugMiddle(handlers.HandleUpdateTodo(logger, queries, validate)))
	mux.Handle("DELETE /todos/{id}", debugMiddle(handlers.HandleDeleteTodo(logger, queries, validate)))
}

func debugMiddleware(l gsdlogger.Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			routeLogger := l.With(
				"method", r.Method,
				"route", r.URL.Path,
			)
			routeLogger.InfoContext(r.Context(), "hitStart")
			now := time.Now()
			h.ServeHTTP(w, r)
			routeLogger.DebugContext(r.Context(), "hitEnd", "timeTaken", time.Since(now))
		})
	}
}
