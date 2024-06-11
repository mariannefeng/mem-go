-- name: GetBook :one
SELECT * FROM book
WHERE id = $1 LIMIT 1;

-- name: GetBooks :many
SELECT * FROM book
ORDER BY created_at DESC;

-- name: GetEntriesByBook :many
SELECT * FROM entry
WHERE book_id = $1
ORDER BY created_at DESC;

-- name: CreateBook :one
INSERT INTO book (
  name
) VALUES (
  $1
)
RETURNING *;

-- name: CreateEntry :one
INSERT INTO entry (
  book_id, type, content, key
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: DeleteEntry :exec
DELETE FROM entry
WHERE id = $1 AND book_id = $2;