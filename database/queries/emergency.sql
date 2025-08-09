-- name: CreateEmergencyContact :one
INSERT INTO emergency (
    resident_id, name, relation, phone, email
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetEmergencyContact :one
SELECT *
FROM emergency
WHERE id = $1
AND deleted_at IS NULL;

-- name: GetEmergencyContactByResident :one
SELECT *
FROM emergency
WHERE resident_id = $1
AND deleted_at IS NULL;

-- name: ListEmergencyContacts :many
SELECT *
FROM emergency
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: UpdateEmergencyContact :one
UPDATE emergency
SET
    name     = COALESCE(sqlc.narg('name'), name),
    relation = COALESCE(sqlc.narg('relation'), relation),
    phone    = COALESCE(sqlc.narg('phone'), phone),
    email    = COALESCE(sqlc.narg('email'), email)
WHERE id = sqlc.arg('id') AND deleted_at IS NULL
RETURNING *;

-- name: DeleteEmergencyContact :exec
UPDATE emergency
SET deleted_at = NOW()
WHERE id = $1
AND deleted_at IS NULL;

-- name: HardDeleteEmergencyContact :exec
DELETE FROM emergency
WHERE id = $1 AND deleted_at IS NOT NULL;
