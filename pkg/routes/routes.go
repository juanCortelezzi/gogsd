package routes

import (
	"net/http"
	"time"

	"github.com/juancortelezzi/gogsd/pkg/gsdlogger"
	"github.com/juancortelezzi/gogsd/pkg/handlers"
)

func AddRoutes(mux *http.ServeMux, l gsdlogger.Logger) {
	debugMiddle := debugMiddleware(l)

	mux.Handle("GET /ping", handlers.HandlePing())
	mux.Handle("GET /hello/{name}", debugMiddle(handlers.HandleHello(l)))
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
			routeLogger.InfoContext(r.Context(), "hitEnd", "timeTaken", time.Since(now))
		})
	}
}
