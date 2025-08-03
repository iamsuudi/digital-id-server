-- name: CreateKebele :one
INSERT INTO kebele (name, lat, lon, city_id, subcity_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateKebele :one
UPDATE kebele
SET name = $2, lat = $3, lon = $4, city_id = $5, subcity_id = $6
WHERE id = $1
RETURNING *;

-- name: GetKebele :one
SELECT k.*, c.name as city_name, sc.name as subcity_name, u.id as executive_id,
    CONCAT_WS(' ', u.first_name, u.second_name, u.last_name) AS executive_name
FROM kebele k
LEFT JOIN city c ON c.id = k.city_id
LEFT JOIN subcity sc ON sc.id = k.subcity_id
LEFT JOIN "user" u ON u.kebele_id = k.id AND u.role_slug = 'executive'
WHERE k.id = $1 AND k.deleted_at IS NULL;

-- name: ListKebeles :many
SELECT k.*, c.name as city_name, sc.name as subcity_name, u.id as executive_id,
    CONCAT_WS(' ', u.first_name, u.second_name, u.last_name) AS executive_name
FROM kebele k
LEFT JOIN city c ON c.id = k.city_id
LEFT JOIN subcity sc ON sc.id = k.subcity_id
LEFT JOIN "user" u ON u.kebele_id = k.id AND u.role_slug = 'executive'
WHERE k.deleted_at IS NULL
ORDER BY k.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountListKebeles :one
SELECT COUNT(*)
FROM kebele
WHERE deleted_at IS NULL;

-- name: SearchKebeles :many
SELECT k.*, c.name as city_name, sc.name as subcity_name, u.id as executive_id,
    CONCAT_WS(' ', u.first_name, u.second_name, u.last_name) AS executive_name,
    similarity(CONCAT_WS(' ', k.name, sc.name, c.name), sqlc.arg('query')) AS sim
FROM kebele k
LEFT JOIN city c ON c.id = k.city_id
LEFT JOIN subcity sc ON sc.id = k.subcity_id
LEFT JOIN "user" u ON u.kebele_id = k.id
WHERE k.deleted_at IS NULL AND
    similarity(CONCAT_WS(' ', k.name, sc.name, c.name), sqlc.arg('query')) > 0.2
ORDER BY sim DESC, k.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountSearchKebeles :one
SELECT COUNT(*)
FROM kebele
WHERE deleted_at IS NULL AND
    similarity(name, sqlc.arg('query')) > 0.2;
