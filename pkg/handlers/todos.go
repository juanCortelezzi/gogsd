package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/juancortelezzi/gogsd/pkg/database"
	"github.com/juancortelezzi/gogsd/pkg/gsdlogger"
)

func HandleListTodos(logger gsdlogger.Logger, queries *database.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		todos, err := queries.ListTodos(r.Context())
		if err != nil {
			http.Error(w, "could not get todos from db", http.StatusInternalServerError)
			return
		}

		todosJson, err := json.Marshal(todos)
		if err != nil {
			http.Error(w, "could not marshal todos", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(todosJson)
	})
}

func HandleCreateTodo(
	logger gsdlogger.Logger,
	queries *database.Queries,
	validate *validator.Validate,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var todoParams struct {
			Description string `validate:"min=1,max=255,ascii"`
			Done        bool
		}

		if err := json.NewDecoder(r.Body).Decode(&todoParams); err != nil {
			http.Error(w, "could not decode todo from body", http.StatusBadRequest)
			return
		}

		if err := validate.Struct(todoParams); err != nil {
			logger.DebugContext(r.Context(), "validation fail", "err", err)
			formattedError := fmt.Errorf("validation fail: %w", err)
			http.Error(w, formattedError.Error(), http.StatusBadRequest)
			return
		}

		logger.DebugContext(r.Context(), "creating todo", "requestParams", todoParams)
		todo, err := queries.CreateTodo(r.Context(), database.CreateTodoParams{
			Description: todoParams.Description,
			Done:        todoParams.Done,
		})

		if err != nil {
			logger.ErrorContext(r.Context(), "could not save todo in database", "err", err)
			http.Error(w, "could not save todo in database", http.StatusInternalServerError)
			return
		}

		todoJson, err := json.Marshal(todo)
		if err != nil {
			logger.ErrorContext(r.Context(), "could not marshal todo", "err", err)
			http.Error(w, "could not marshal todo", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		w.Write(todoJson)
	})
}

func HandleUpdateTodo(
	logger gsdlogger.Logger,
	queries *database.Queries,
	validate *validator.Validate,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := r.PathValue("id")
		if idParam == "" {
			logger.ErrorContext(r.Context(), "could not find id in path")
			http.Error(w, "could not find id in path", http.StatusInternalServerError)
			return
		}

		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			logger.DebugContext(r.Context(), "could not parse id", "err", err)
			http.Error(w, "could not parse id", http.StatusBadRequest)
			return
		}

		var todoParams struct {
			Description string `validate:"min=1,max=255,ascii"`
			Done        bool
		}

		if err := json.NewDecoder(r.Body).Decode(&todoParams); err != nil {
			logger.DebugContext(r.Context(), "could not decode todo from body", "err", err)
			http.Error(w, "could not decode todo from body", http.StatusBadRequest)
			return
		}

		if err := validate.Struct(todoParams); err != nil {
			logger.DebugContext(r.Context(), "validation fail", "err", err)
			formattedError := fmt.Errorf("validation fail: %w", err)
			http.Error(w, formattedError.Error(), http.StatusBadRequest)
			return
		}

		logger.DebugContext(r.Context(), "updating todo", "requestParams", todoParams)
		todo, err := queries.UpdateTodo(r.Context(), database.UpdateTodoParams{
			Description: todoParams.Description,
			Done:        todoParams.Done,
			ID:          id,
		})

		if err != nil {
			logger.ErrorContext(r.Context(), "could not update todo in database", "err", err)
			http.Error(w, "could not update todo in database", http.StatusInternalServerError)
			return
		}

		todoJson, err := json.Marshal(todo)
		if err != nil {
			logger.ErrorContext(r.Context(), "could not marshal todo", "err", err)
			http.Error(w, "could not marshal todo", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(todoJson)
	})
}

func HandleDeleteTodo(
	logger gsdlogger.Logger,
	queries *database.Queries,
	validate *validator.Validate,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := r.PathValue("id")
		if idParam == "" {
			logger.ErrorContext(r.Context(), "could not find id in path")
			http.Error(w, "could not find id in path", http.StatusInternalServerError)
			return
		}

		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			logger.DebugContext(r.Context(), "could not parse id", "err", err)
			http.Error(w, "could not parse id", http.StatusBadRequest)
			return
		}

		logger.DebugContext(r.Context(), "deleting todo", "id", id)
		if err := queries.DeleteTodo(r.Context(), id); err != nil {
			logger.ErrorContext(r.Context(), "could not delete todo", "err", err)
			http.Error(w, "could not delete todo in database", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
