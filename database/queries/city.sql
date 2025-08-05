-- name: CreateCity :one
INSERT INTO city (name, lat, lon)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateCity :exec
UPDATE city
SET name = $2, lat = $3, lon = $4
WHERE id = $1;

-- name: GetCity :one
SELECT c.*, u.id as admin_id, 
    CONCAT_WS(' ', u.first_name, u.second_name, u.last_name) AS admin_name
FROM city c
LEFT JOIN "user" u ON u.city_id = c.id AND u.role_slug = 'admin'
WHERE c.id = $1 AND c.deleted_at IS NULL;

-- name: SoftDeleteCity :exec
UPDATE city
SET deleted_at = NOW()
WHERE id = $1;

-- name: GetSubCitiesForCity :many
SELECT s.*, c.name as city_name, u.id as manager_id, 
    CONCAT_WS(' ', u.first_name, u.second_name, u.last_name) AS manager_name
FROM subcity s
JOIN city c ON c.id = s.city_id
LEFT JOIN "user" u ON u.subcity_id = s.id AND u.role_slug = 'manager'
WHERE s.deleted_at IS NULL AND c.id = $1
ORDER BY s.created_at DESC;

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

-- name: CountSearchCities :one
SELECT COUNT(*)
FROM city
WHERE deleted_at IS NULL AND similarity(name, sqlc.arg('query')) > 0.2;
