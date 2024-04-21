package tests

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/juancortelezzi/gogsd/pkg/database"
	"github.com/juancortelezzi/gogsd/pkg/gsdlogger"
	"github.com/juancortelezzi/gogsd/pkg/server"
)

func TestHelloRoute(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	{
		logger := gsdlogger.NewLogger(os.Stdout, slog.LevelDebug)
		go server.Run(ctx, logger, testLookupEnv)
		err := waitForReady(ctx, logger, getBaseUrl()+"/ping")
		if err != nil {
			t.Fatal(err)
		}
	}

	resp, err := http.Get(getBaseUrl() + "/hello/world")
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	data, err := bufio.NewReaderSize(resp.Body, 1024).ReadBytes('\n')
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(data, []byte("Hello, world!\n")) {
		t.Fatalf("expected Hello, world! got %s\n", data)
	}
}

func TestListTodosRoute(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	{
		logger := gsdlogger.NewLogger(os.Stdout, slog.LevelDebug)
		go server.Run(ctx, logger, testLookupEnv)
		err := waitForReady(ctx, logger, getBaseUrl()+"/ping")
		if err != nil {
			t.Fatal(err)
		}
	}

	resp, err := http.Get(getBaseUrl() + "/todos")
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf(
			"expected status code to be %d but got %d",
			http.StatusOK,
			resp.StatusCode,
		)
	}

	var todos []database.Todo
	if err := json.NewDecoder(resp.Body).Decode(&todos); err != nil {
		t.Fatal(err)
	}

	if len(todos) != 0 {
		t.Fatalf("expected 0 todos got %d\n", len(todos))
	}
}

func TestCreateTodoRoute(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	{
		logger := gsdlogger.NewLogger(os.Stdout, slog.LevelDebug)
		go server.Run(ctx, logger, testLookupEnv)
		err := waitForReady(ctx, logger, getBaseUrl()+"/ping")
		if err != nil {
			t.Fatal(err)
		}
	}

	todoParams := `{ "description": "finish this server", "done": true }`

	resp, err := http.Post(
		getBaseUrl()+"/todos",
		"application/json",
		strings.NewReader(todoParams),
	)

	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status code to be %d but got %d", http.StatusCreated, resp.StatusCode)
	}

	var todo database.Todo
	if err := json.NewDecoder(resp.Body).Decode(&todo); err != nil {
		t.Fatal(err)
	}

	expected := "finish this server"
	if todo.Description != "finish this server" {
		t.Fatalf("expected description to be '%s' but got '%s'", expected, todo.Description)
	}

	if !todo.Done {
		t.Fatalf("expected done to be true but got false")
	}
}

func TestCreateTodoRouteFail(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	{
		logger := gsdlogger.NewLogger(os.Stdout, slog.LevelDebug)
		go server.Run(ctx, logger, testLookupEnv)
		err := waitForReady(ctx, logger, getBaseUrl()+"/ping")
		if err != nil {
			t.Fatal(err)
		}
	}

	todoParams := `{ "done": true }`

	resp, err := http.Post(
		getBaseUrl()+"/todos",
		"application/json",
		strings.NewReader(todoParams),
	)

	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status code to be %d but got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestUpdateTodoRoute(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	{
		logger := gsdlogger.NewLogger(os.Stdout, slog.LevelDebug)
		go server.Run(ctx, logger, testLookupEnv)
		err := waitForReady(ctx, logger, getBaseUrl()+"/ping")
		if err != nil {
			t.Fatal(err)
		}
	}

	todoParams := `{ "description": "finish this server", "done": true }`

	resp, err := http.Post(
		getBaseUrl()+"/todos",
		"application/json",
		strings.NewReader(todoParams),
	)

	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status code to be %d but got %d", http.StatusCreated, resp.StatusCode)
	}

	var todo database.Todo
	if err := json.NewDecoder(resp.Body).Decode(&todo); err != nil {
		t.Fatal(err)
	}

	client := &http.Client{}
	updateTodoParams := `{ "description": "finish this test", "done": false }`

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		fmt.Sprintf("%s/todos/%d", getBaseUrl(), todo.ID),
		strings.NewReader(updateTodoParams),
	)

	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code to be %d but got %d", http.StatusOK, resp.StatusCode)
	}

	var updatedTodo database.Todo
	if err := json.NewDecoder(resp.Body).Decode(&todo); err != nil {
		t.Fatal(err)
	}

	if updatedTodo.ID == todo.ID {
		t.Fatalf("expected id to be %d but got %d", todo.ID, updatedTodo.ID)
	}

	if updatedTodo.Done {
		t.Fatalf("expected done to be false but got %t", updatedTodo.Done)
	}

	if updatedTodo.Description == "finish this test" {
		t.Fatalf(`expected description to be "finish this test" but got %s`, updatedTodo.Description)
	}
}

func TestDeleteTodoRoute(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	{
		logger := gsdlogger.NewLogger(os.Stdout, slog.LevelDebug)
		go server.Run(ctx, logger, testLookupEnv)
		err := waitForReady(ctx, logger, getBaseUrl()+"/ping")
		if err != nil {
			t.Fatal(err)
		}
	}

	todoParams := `{ "description": "finish this server", "done": true }`

	resp, err := http.Post(
		getBaseUrl()+"/todos",
		"application/json",
		strings.NewReader(todoParams),
	)

	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status code to be %d but got %d", http.StatusCreated, resp.StatusCode)
	}

	var todo database.Todo
	if err := json.NewDecoder(resp.Body).Decode(&todo); err != nil {
		t.Fatal(err)
	}

	client := &http.Client{}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("%s/todos/%d", getBaseUrl(), todo.ID),
		nil,
	)

	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code to be %d but got %d", http.StatusOK, resp.StatusCode)
	}
}
