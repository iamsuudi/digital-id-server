-- name: CreateCity :one
INSERT INTO city (name)
VALUES ($1)
RETURNING *;

-- name: GetCity :one
SELECT c.*,
a.first_name,
a.second_name,
a.last_name
FROM city c
LEFT JOIN actor a ON a.id = c.admin_id
WHERE c.id = $1 AND c.deleted_at IS NULL;

-- name: ListCities :many
SELECT c.*,
a.first_name,
a.second_name,
a.last_name,
COUNT(*) OVER() as count
FROM city c
LEFT JOIN actor a ON a.id = c.admin_id
WHERE c.deleted_at IS NULL
ORDER BY c.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CreateKebele :one
INSERT INTO kebele (name, city_id)
VALUES ($1, $2)
RETURNING *;