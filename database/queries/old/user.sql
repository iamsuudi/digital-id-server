-- name: CreateUser :one
INSERT INTO actor (
  first_name, second_name, last_name,
  email, phone, password_hash, role
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: ListUsers :many
SELECT *, CONCAT_WS(' ', first_name, second_name, last_name) AS full_name
FROM actor
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountListUsers :one
SELECT COUNT(*)
FROM actor
WHERE deleted_at IS NULL;

-- name: GetUserByEmail :one
SELECT *, CONCAT_WS(' ', first_name, second_name, last_name) AS full_name
FROM actor
WHERE email = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: GetUserByID :one
SELECT *, CONCAT_WS(' ', first_name, second_name, last_name) AS full_name
FROM actor
WHERE id = $1 AND deleted_at IS NULL;

-- name: SoftDeleteUser :exec
UPDATE actor
SET deleted_at = NOW()
WHERE id = $1;

-- name: SearchUsers :many
SELECT *, CONCAT_WS(' ', first_name, second_name, last_name) AS full_name
FROM actor
WHERE deleted_at IS NULL
AND search_vector @@ to_tsquery('english', sqlc.arg('query') || ':*')
ORDER BY
ts_rank(search_vector, to_tsquery('english', sqlc.arg('query') || ':*')) DESC,
similarity(first_name || ' ' || second_name || ' ' || last_name, sqlc.arg('query')) DESC,
created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountUsersSearch :one
SELECT COUNT(*)
FROM actor
WHERE deleted_at IS NULL
AND search_vector @@ plainto_tsquery('english', sqlc.arg('query'));
