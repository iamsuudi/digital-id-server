-- name: CreateIDCard :one
INSERT INTO idcard (
    resident_id, number, issue_date, expiry_date, issue_place
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetIDCard :one
SELECT *
FROM idcard
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetIDCardByResident :one
SELECT *
FROM idcard
WHERE resident_id = $1 AND deleted_at IS NULL;

-- name: GetIDCardByNumber :one
SELECT *
FROM idcard
WHERE number = $1 AND deleted_at IS NULL;

-- name: ListIDCards :many
SELECT *
FROM idcard
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: UpdateIDCard :one
UPDATE idcard
SET
    number       = COALESCE(sqlc.narg('number'), number),
    issue_date   = COALESCE(sqlc.narg('issue_date'), issue_date),
    expiry_date  = COALESCE(sqlc.narg('expiry_date'), expiry_date),
    issue_place  = COALESCE(sqlc.narg('issue_place'), issue_place)
WHERE id = sqlc.arg('id') AND deleted_at IS NULL
RETURNING *;

-- name: DeleteIDCard :exec
UPDATE idcard
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: HardDeleteIDCard :exec
DELETE FROM idcard
WHERE id = $1 AND deleted_at IS NOT NULL;
