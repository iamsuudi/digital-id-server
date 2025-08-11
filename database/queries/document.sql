-- name: CreateDocument :one
INSERT INTO document (
    resident_id, url, status
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetDocument :one
SELECT *
FROM document
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetDocumentByResident :one
SELECT *
FROM document
WHERE resident_id = $1 AND deleted_at IS NULL;

-- name: ListDocuments :many
SELECT *
FROM document
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: UpdateDocument :one
UPDATE document
SET
    url    = COALESCE(sqlc.narg('url'), url),
    status = COALESCE(sqlc.narg('status'), status)
WHERE id = sqlc.arg('id') AND deleted_at IS NULL
RETURNING *;

-- name: DeleteDocument :exec
UPDATE document
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: HardDeleteDocument :exec
DELETE FROM document
WHERE id = $1 AND deleted_at IS NOT NULL;

-- name: GetResidentDocuments :many
SELECT *
FROM document
WHERE resident_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: DeleteResidentDocuments :exec
DELETE FROM document
WHERE resident_id = $1 AND deleted_at IS NULL;
