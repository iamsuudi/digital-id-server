-- name: CreateResident :one
INSERT INTO resident (
  email, first_name, second_name, last_name, birth_date, gender, phone, address_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING id;

-- name: UpdateResident :exec
UPDATE resident SET
  email = $2, first_name = $3, second_name = $4, last_name = $5,
  birth_date = $6, gender = $7, phone = $8
WHERE id = $1;

-- name: UpdateResidentAddress :exec
UPDATE resident SET address_id = $1 WHERE id = $2;

-- name: GetResident :one
SELECT sqlc.embed(resident)
    -- sqlc.embed(address), sqlc.embed(biometric), sqlc.embed(document),
    -- sqlc.embed(employment), sqlc.embed(emergency), sqlc.embed(idcard)
FROM resident
-- LEFT JOIN address    ON resident.address_id = address.id
-- LEFT JOIN biometric  ON resident.id = biometric.resident_id
-- LEFT JOIN document   ON resident.id = document.resident_id
-- LEFT JOIN employment ON resident.id = employment.resident_id
-- LEFT JOIN emergency  ON resident.id = emergency.resident_id
-- LEFT JOIN idcard     ON resident.id = idcard.resident_id
WHERE resident.id = $1;

-- name: DeleteResident :exec
DELETE FROM resident WHERE id = $1;

-- name: DeleteAllResidents :exec
DELETE FROM resident;

-- name: ListResidents :many
SELECT *, CONCAT_WS(' ', first_name, second_name, last_name) AS full_name
FROM resident
WHERE deleted_at IS NULL
ORDER BY created_at ASC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountListResidents :one
SELECT COUNT(*)
FROM resident
WHERE deleted_at IS NULL;

-- name: SearchResidents :many
SELECT *, CONCAT_WS(' ', first_name, second_name, last_name) AS full_name,
    similarity(CONCAT_WS(' ', first_name, second_name, last_name), sqlc.arg('query')) AS sim
FROM resident
WHERE deleted_at IS NULL AND
    similarity(CONCAT_WS(' ', first_name, second_name, last_name), sqlc.arg('query')) > 0.2
ORDER BY sim DESC, created_at ASC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountSearchResidents :one
SELECT COUNT(*)
FROM resident
WHERE deleted_at IS NULL AND
    similarity(CONCAT_WS(' ', first_name, second_name, last_name), sqlc.arg('query')) > 0.2;
