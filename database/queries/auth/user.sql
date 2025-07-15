-- name: CreateUser :one
INSERT INTO actor (
  first_name, second_name, last_name,
  email, phone, password_hash, role
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING id;

-- name: GetUserByEmail :one
SELECT * FROM actor
WHERE email = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM actor
WHERE id = $1 AND deleted_at IS NULL;

-- name: SoftDeleteUser :exec
UPDATE actor
SET deleted_at = NOW()
WHERE id = $1;
