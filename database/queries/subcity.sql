-- name: CreateSubCity :one
INSERT INTO subcity (name, city_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetSubCity :one
SELECT s.*, u.id as manager_id, CONCAT_WS(' ', u.first_name, u.second_name, u.last_name) AS manager_name
FROM subcity s
LEFT JOIN "user" u ON u.subcity_id = s.id
WHERE s.id = $1 AND s.deleted_at IS NULL;

-- name: ListSubCities :many
SELECT s.*, u.id as manager_id, CONCAT_WS(' ', u.first_name, u.second_name, u.last_name) AS manager_name
FROM subcity s
LEFT JOIN "user" u ON u.subcity_id = s.id
WHERE s.deleted_at IS NULL
ORDER BY s.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountListSubcities :one
SELECT COUNT(*)
FROM subcity
WHERE deleted_at IS NULL;
