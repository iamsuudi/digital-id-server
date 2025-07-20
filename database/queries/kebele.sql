-- name: CreateKebele :one
INSERT INTO kebele (name, city_id)
VALUES ($1, $2)
RETURNING *;