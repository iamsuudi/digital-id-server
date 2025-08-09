-- name: CreateAddress :one
INSERT INTO address (
    house_number, kebele_id, subcity_id, city_id
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetAddress :one
SELECT *
FROM address
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetAddressByLocations :one
SELECT *
FROM address
WHERE city_id = $1 AND subcity_id = $2 AND kebele_id = $3
    AND house_number = $4 AND deleted_at IS NULL;

-- name: ListAddresses :many
SELECT *
FROM address
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: UpdateAddress :one
UPDATE address
SET
    house_number = COALESCE(sqlc.narg('house_number'), house_number),
    kebele_id    = COALESCE(sqlc.narg('kebele_id'), kebele_id),
    subcity_id   = COALESCE(sqlc.narg('subcity_id'), subcity_id),
    city_id      = COALESCE(sqlc.narg('city_id'), city_id)
WHERE id = sqlc.arg('id')
AND deleted_at IS NULL
RETURNING *;

-- name: DeleteAddress :exec
UPDATE address
SET deleted_at = NOW()
WHERE id = $1
AND deleted_at IS NULL;

-- name: HardDeleteAddress :exec
DELETE FROM address
WHERE id = $1
AND deleted_at IS NOT NULL;

-- name: GetRandomLocation :one
SELECT 
    c.id   AS city_id,
    s.id   AS subcity_id,
    k.id   AS kebele_id
FROM kebele k
JOIN subcity s ON k.subcity_id = s.id
JOIN city c ON s.city_id = c.id
ORDER BY RANDOM()
LIMIT 1;

