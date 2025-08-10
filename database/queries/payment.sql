-- name: CreatePayment :one
INSERT INTO payment (
    resident_id, status
) VALUES (
    $1, $2
)
RETURNING *;

-- name: GetPayment :one
SELECT *
FROM payment
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetPaymentByResident :one
SELECT *
FROM payment
WHERE resident_id = $1 AND deleted_at IS NULL;

-- name: ListPayments :many
SELECT *
FROM payment
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: UpdatePayment :one
UPDATE payment
SET
    amount      = COALESCE(sqlc.narg('amount'), amount),
    description = COALESCE(sqlc.narg('description'), description),
    status      = COALESCE(sqlc.narg('status'), status),
    reference   = COALESCE(sqlc.narg('reference'), reference),
    method      = COALESCE(sqlc.narg('method'), method)
WHERE id = sqlc.arg('id') AND deleted_at IS NULL
RETURNING *;

-- name: DeletePayment :exec
UPDATE payment
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: HardDeletePayment :exec
DELETE FROM payment
WHERE id = $1 AND deleted_at IS NOT NULL;
