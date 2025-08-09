-- name: CreateAdditional :one
INSERT INTO additional (
    resident_id,  national_id,  marital_status,  religion,  ethnicity,  disability,  education_level,  languages_spoken
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetAdditional :one
SELECT *
FROM additional
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetAdditionalByResident :one
SELECT *
FROM additional
WHERE resident_id = $1 AND deleted_at IS NULL;

-- name: ListAdditionals :many
SELECT *
FROM additional
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: UpdateAdditional :one
UPDATE additional
SET
    national_id      = COALESCE(sqlc.narg('national_id'), national_id),
    marital_status   = COALESCE(sqlc.narg('marital_status'), marital_status),
    religion         = COALESCE(sqlc.narg('religion'), religion),
    ethnicity        = COALESCE(sqlc.narg('ethnicity'), ethnicity),
    disability       = COALESCE(sqlc.narg('disability'), disability),
    education_level  = COALESCE(sqlc.narg('education_level'), education_level),
    languages_spoken = COALESCE(sqlc.narg('languages_spoken'), languages_spoken)
WHERE id = sqlc.arg('id')
AND deleted_at IS NULL
RETURNING *;

-- name: DeleteAdditional :exec
UPDATE additional
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: HardDeleteAdditional :exec
DELETE FROM additional
WHERE id = $1 AND deleted_at IS NOT NULL;
