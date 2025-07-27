-- name: CreateKebele :one
INSERT INTO kebele (name, city_id, subcity_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateKebele :one
UPDATE kebele
SET name = $2, city_id = $3, subcity_id = $4
WHERE id = $1
RETURNING *;

-- name: GetKebele :one
SELECT k.*, c.name as city_name, sc.name as subcity_name, u.id as executive_id, 
    CONCAT_WS(' ', u.first_name, u.second_name, u.last_name) AS executive_name
FROM kebele k
LEFT JOIN city c ON c.id = k.city_id
LEFT JOIN subcity sc ON sc.id = k.subcity_id
LEFT JOIN "user" u ON u.kebele_id = k.id
WHERE k.id = $1 AND k.deleted_at IS NULL;

-- name: ListKebeles :many
SELECT k.*, c.name as city_name, sc.name as subcity_name, u.id as executive_id,
    CONCAT_WS(' ', u.first_name, u.second_name, u.last_name) AS executive_name
FROM kebele k
LEFT JOIN city c ON c.id = k.city_id
LEFT JOIN subcity sc ON sc.id = k.subcity_id
LEFT JOIN "user" u ON u.kebele_id = k.id
WHERE k.deleted_at IS NULL
ORDER BY k.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountListKebeles :one
SELECT COUNT(*)
FROM kebele
WHERE deleted_at IS NULL;
