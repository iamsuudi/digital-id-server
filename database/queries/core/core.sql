-- name: CreateRegion :one
INSERT INTO region (id, name, created_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateCity :one
INSERT INTO city (id, name, region_id, created_at)
VALUES ($1, $2, $3, $4)
RETURNING *;