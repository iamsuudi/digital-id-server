-- name: CreateCity :one
INSERT INTO city (name)
VALUES ($1)
RETURNING *;

-- name: UpdateCity :one
UPDATE city
SET name = $2
WHERE id = $1
RETURNING *;

-- name: GetCity :one
SELECT c.*, u.id as admin_id, 
    CONCAT_WS(' ', u.first_name, u.second_name, u.last_name) AS admin_name
FROM city c
LEFT JOIN "user" u ON u.city_id = c.id AND u.role_slug = 'admin'
WHERE c.id = $1 AND c.deleted_at IS NULL;

-- name: ListCities :many
SELECT c.*, u.id as admin_id, 
    CONCAT_WS(' ', u.first_name, u.second_name, u.last_name) AS admin_name
FROM city c
LEFT JOIN "user" u ON u.city_id = c.id AND u.role_slug = 'admin'
WHERE c.deleted_at IS NULL
ORDER BY c.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountListCities :one
SELECT COUNT(*)
FROM city
WHERE deleted_at IS NULL;

-- name: SearchCities :many
SELECT *, u.id as admin_id, similarity(c.name, sqlc.arg('query')) AS sim,
    CONCAT_WS(' ', u.first_name, u.second_name, u.last_name) AS admin_name
FROM city c
LEFT JOIN "user" u ON u.city_id = c.id AND u.role_slug = 'admin'
WHERE c.deleted_at IS NULL AND similarity(c.name, sqlc.arg('query')) > 0.2
ORDER BY sim DESC, c.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountCitiesSearch :one
SELECT COUNT(*)
FROM city
WHERE deleted_at IS NULL AND similarity(name, sqlc.arg('query')) > 0.2;
