-- name: GetTodo :one
SELECT * FROM todos
WHERE id = ? LIMIT 1;

-- name: ListTodos :many
SELECT * FROM todos
ORDER BY created_at;

-- name: CreateTodo :one
INSERT INTO todos (
  description, 
  done
) VALUES (
  ?, ?
)
RETURNING *;

-- name: UpdateTodo :one
UPDATE todos
set description = ?,
done = ?
WHERE id = ?
RETURNING *;

-- name: DeleteTodo :exec
DELETE FROM todos
WHERE id = ?;
