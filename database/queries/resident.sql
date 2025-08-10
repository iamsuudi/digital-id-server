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
SELECT sqlc.embed(resident), sqlc.embed(address), sqlc.embed(biometric),
    sqlc.embed(additional), sqlc.embed(employment), sqlc.embed(emergency)
FROM resident
JOIN address    ON resident.address_id = address.id
JOIN biometric  ON resident.id = biometric.resident_id
JOIN document   ON resident.id = document.resident_id
JOIN employment ON resident.id = employment.resident_id
JOIN emergency  ON resident.id = emergency.resident_id
JOIN additional ON resident.id = additional.resident_id
WHERE resident.id = $1;

-- name: DeleteResident :exec
DELETE FROM resident WHERE id = $1;

-- name: DeleteAllResidents :exec
DELETE FROM resident;

-- name: ListResidents :many
SELECT r.*, b.face_url, CONCAT_WS(' ', r.first_name, r.second_name, r.last_name) AS full_name
FROM resident r
JOIN biometric b   ON r.id = b.resident_id
JOIN payment p     ON p.resident_id = r.id
WHERE r.deleted_at IS NULL AND p.status = 'verified'
ORDER BY r.created_at ASC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountListResidents :one
SELECT COUNT(*)
FROM resident r
JOIN payment p     ON p.resident_id = r.id
WHERE r.deleted_at IS NULL AND p.status = 'verified';

-- name: SearchResidents :many
SELECT r.*, b.face_url, CONCAT_WS(' ', r.first_name, r.second_name, r.last_name) AS full_name,
    similarity(CONCAT_WS(' ', r.first_name, r.second_name, r.last_name), sqlc.arg('query')) AS sim
FROM resident r
JOIN biometric b   ON r.id = b.resident_id
JOIN payment p     ON p.resident_id = r.id
WHERE r.deleted_at IS NULL AND p.status = 'verified' AND
    similarity(CONCAT_WS(' ', r.first_name, r.second_name, r.last_name), sqlc.arg('query')) > 0.2
ORDER BY sim DESC, r.created_at ASC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountSearchResidents :one
SELECT COUNT(*)
FROM resident r
JOIN payment p     ON p.resident_id = r.id
WHERE r.deleted_at IS NULL AND p.status = 'verified' AND
    similarity(CONCAT_WS(' ', r.first_name, r.second_name, r.last_name), sqlc.arg('query')) > 0.2;

-- name: ListUnpaidResidents :many
SELECT sqlc.embed(r), sqlc.embed(p), b.face_url, CONCAT_WS(' ', r.first_name, r.second_name, r.last_name) AS full_name
FROM resident r
JOIN biometric b   ON r.id = b.resident_id
JOIN payment p     ON p.resident_id = r.id
WHERE r.deleted_at IS NULL AND p.status != 'verified'
ORDER BY r.created_at ASC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountListUnpaidResidents :one
SELECT COUNT(*)
FROM resident r
JOIN payment p     ON p.resident_id = r.id
WHERE r.deleted_at IS NULL AND p.status != 'verified';

-- name: SearchUnpaidResidents :many
SELECT sqlc.embed(r), sqlc.embed(p), b.face_url, CONCAT_WS(' ', r.first_name, r.second_name, r.last_name) AS full_name,
    similarity(CONCAT_WS(' ', r.first_name, r.second_name, r.last_name), sqlc.arg('query')) AS sim
FROM resident r
JOIN biometric b   ON r.id = b.resident_id
JOIN payment p     ON p.resident_id = r.id
WHERE r.deleted_at IS NULL AND p.status != 'verified' AND
    similarity(CONCAT_WS(' ', r.first_name, r.second_name, r.last_name), sqlc.arg('query')) > 0.2
ORDER BY r.created_at ASC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountSearchUnpaidResidents :one
SELECT COUNT(*)
FROM resident r
JOIN payment p     ON p.resident_id = r.id
WHERE r.deleted_at IS NULL AND p.status != 'verified' AND
    similarity(CONCAT_WS(' ', r.first_name, r.second_name, r.last_name), sqlc.arg('query')) > 0.2;
