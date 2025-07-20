-- name: CreateCity :one
INSERT INTO city (name)
VALUES ($1)
RETURNING *;

-- name: GetCity :one
SELECT c.*, u.id as admin_id, CONCAT_WS(' ', u.first_name, u.second_name, u.last_name) AS admin_name
FROM city c
LEFT JOIN "user" u ON u.city_id = c.id
WHERE c.id = $1 AND c.deleted_at IS NULL;

-- name: ListCities :many
SELECT c.*, u.id as admin_id, CONCAT_WS(' ', u.first_name, u.second_name, u.last_name) AS admin_name
FROM city c
LEFT JOIN "user" u ON u.city_id = c.id
WHERE c.deleted_at IS NULL
ORDER BY c.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountListCities :one
SELECT COUNT(*)
FROM city
WHERE deleted_at IS NULL;
