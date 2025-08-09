-- name: CreateEmployment :one
INSERT INTO employment (
    resident_id, status, type, occupation, employer_name, work_address
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetEmployment :one
SELECT *
FROM employment
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetEmploymentByResident :one
SELECT *
FROM employment
WHERE resident_id = $1 AND deleted_at IS NULL;

-- name: ListEmployments :many
SELECT *
FROM employment
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: UpdateEmployment :one
UPDATE employment
SET
    status        = COALESCE(sqlc.narg('status'), status),
    type          = COALESCE(sqlc.narg('type'), type),
    occupation    = COALESCE(sqlc.narg('occupation'), occupation),
    employer_name = COALESCE(sqlc.narg('employer_name'), employer_name),
    work_address  = COALESCE(sqlc.narg('work_address'), work_address)
WHERE id = sqlc.arg('id') AND deleted_at IS NULL
RETURNING *;

-- name: DeleteEmployment :exec
UPDATE employment
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: HardDeleteEmployment :exec
DELETE FROM employment
WHERE id = $1 AND deleted_at IS NOT NULL;
