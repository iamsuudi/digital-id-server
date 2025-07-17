-- name: GetAllResidents :many
SELECT *
FROM resident
WHERE deleted_at IS NULL
ORDER BY created_at DESC;
