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
	logMiddle := logMiddleware(logger)

	mux.Handle("GET /ping", handlers.HandlePing())
	mux.Handle("GET /hello/{name}", handlers.HandleHello(logger))

	mux.Handle("GET /todos", logMiddle(func(l gsdlogger.Logger) http.Handler {
		return handlers.HandleListTodos(l, queries)
	}))

	mux.Handle("POST /todos", logMiddle(func(l gsdlogger.Logger) http.Handler {
		return handlers.HandleCreateTodo(l, queries, validate)
	}))

	mux.Handle("PUT /todos/{id}", logMiddle(func(l gsdlogger.Logger) http.Handler {
		return handlers.HandleUpdateTodo(l, queries, validate)
	}))

	mux.Handle("DELETE /todos/{id}", logMiddle(func(l gsdlogger.Logger) http.Handler {
		return handlers.HandleDeleteTodo(l, queries, validate)
	}))
}

func logMiddleware(logger gsdlogger.Logger) func(wrapper func(l gsdlogger.Logger) http.Handler) http.Handler {
	return func(wrapper func(l gsdlogger.Logger) http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			l := logger.With("method", r.Method, "path", r.URL.EscapedPath())
			now := time.Now()
			rw := gsdlogger.NewLoggerResponseWritter(w)
			wrapper(l).ServeHTTP(rw, r)
			l.InfoContext(
				r.Context(), "hit",
				"duration", time.Since(now),
				"status", rw.Status(),
			)
		})
	}
}
