-- name: CreateKebele :one
INSERT INTO kebele (name, city_id, subcity_id)
VALUES ($1, $2, $3)
RETURNING *;