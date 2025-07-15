-- name: CreateCity :one
INSERT INTO city (id, name, created_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateKebele :one
INSERT INTO kebele (id, name, city_id, created_at)
VALUES ($1, $2, $3, $4)
RETURNING *;