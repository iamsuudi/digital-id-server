-- name: CreateSubCity :one
INSERT INTO subcity (name, city_id)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateSubCity :one
UPDATE subcity
SET name = $2, city_id = $3
WHERE id = $1
RETURNING *;

-- name: GetSubCity :one
SELECT s.*, c.name as city_name, u.id as manager_id, CONCAT_WS(' ', u.first_name, u.second_name, u.last_name) AS manager_name
FROM subcity s
LEFT JOIN city c ON c.id = s.city_id
LEFT JOIN "user" u ON u.subcity_id = s.id
WHERE s.id = $1 AND s.deleted_at IS NULL;

-- name: ListSubCities :many
SELECT s.*, c.name as city_name, u.id as manager_id, CONCAT_WS(' ', u.first_name, u.second_name, u.last_name) AS manager_name
FROM subcity s
LEFT JOIN city c ON c.id = s.city_id
LEFT JOIN "user" u ON u.subcity_id = s.id
WHERE s.deleted_at IS NULL
ORDER BY s.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountListSubcities :one
SELECT COUNT(*)
FROM subcity
WHERE deleted_at IS NULL;

-- name: SearchSubCities :many
SELECT *, c.name as city_name, u.id as manager_id, CONCAT_WS(' ', u.first_name, u.second_name, u.last_name) AS manager_name,
    similarity(CONCAT_WS(' ', sb.name), sqlc.arg('query')) AS sim
FROM subcity sb
JOIN city c ON c.id = sb.city_id
LEFT JOIN "user" u ON u.subcity_id = sb.id
WHERE sb.deleted_at IS NULL AND
    similarity(CONCAT_WS(' ', sb.name), sqlc.arg('query')) > 0.2
ORDER BY sim DESC, sb.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountSubCitiesSearch :one
SELECT COUNT(*)
FROM subcity
WHERE deleted_at IS NULL AND
    similarity(CONCAT_WS(' ', name), sqlc.arg('query')) > 0.2;

