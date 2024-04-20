package handlers

import (
	"fmt"
	"net/http"

	"github.com/juancortelezzi/gogsd/pkg/gsdlogger"
)

func HandlePing() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Pong")
	})
}

func HandleHello(l gsdlogger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")

		if name == "" {
			l.ErrorContext(r.Context(), "error getting name from url", "name", name)
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, "Hello, %s!\n", name)
	})
}
