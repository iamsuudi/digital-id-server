-- name: CreateResident :one
INSERT INTO resident (
  email, first_name, second_name, last_name, birth_date, gender, phone, address_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: UpdateResident :one
UPDATE resident SET
  email = $2, first_name = $3, second_name = $4, last_name = $5,
  birth_date = $6, gender = $7, phone = $8
WHERE id = $1
RETURNING *;

-- name: UpdateResidentAddress :exec
UPDATE resident SET address_id = $1 WHERE id = $2;

-- name: GetResidentAddress :one
SELECT sqlc.embed(a), sqlc.embed(c), sqlc.embed(s), sqlc.embed(k)
FROM resident
JOIN address a ON resident.address_id = a.id
JOIN city c ON a.city_id = c.id
JOIN subcity s ON a.subcity_id = s.id
JOIN kebele k ON a.kebele_id = k.id
WHERE resident.id = $1 AND resident.deleted_at IS NULL;

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

-- name: GetVerifiedResident :one
SELECT sqlc.embed(resident), sqlc.embed(address), sqlc.embed(biometric),
    sqlc.embed(city), sqlc.embed(subcity), sqlc.embed(kebele)
FROM resident
JOIN address    ON resident.address_id = address.id
JOIN biometric  ON resident.id = biometric.resident_id
JOIN document   ON resident.id = document.resident_id
JOIN city       ON address.city_id = city.id
JOIN subcity    ON address.subcity_id = subcity.id
JOIN kebele     ON address.kebele_id = kebele.id
WHERE resident.id = $1 AND EXISTS (
    SELECT 1 FROM document d
    WHERE d.resident_id = resident.id AND d.status = 'verified'
);

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
AND EXISTS (
    SELECT 1 FROM document d
    WHERE d.resident_id = r.id AND d.status = 'verified'
)
ORDER BY r.created_at ASC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountListResidents :one
SELECT COUNT(*)
FROM resident r
JOIN payment p     ON p.resident_id = r.id
WHERE r.deleted_at IS NULL AND p.status = 'verified'
AND EXISTS (
    SELECT 1 FROM document d
    WHERE d.resident_id = r.id AND d.status = 'verified'
);

-- name: SearchResidents :many
SELECT r.*, b.face_url, CONCAT_WS(' ', r.first_name, r.second_name, r.last_name) AS full_name,
    similarity(CONCAT_WS(' ', r.first_name, r.second_name, r.last_name), sqlc.arg('query')) AS sim
FROM resident r
JOIN biometric b   ON r.id = b.resident_id
JOIN payment p     ON p.resident_id = r.id
WHERE r.deleted_at IS NULL AND p.status = 'verified'
AND EXISTS (
    SELECT 1 FROM document d
    WHERE d.resident_id = r.id AND d.status = 'verified'
) AND similarity(CONCAT_WS(' ', r.first_name, r.second_name, r.last_name), sqlc.arg('query')) > 0.2
ORDER BY sim DESC, r.created_at ASC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountSearchResidents :one
SELECT COUNT(*)
FROM resident r
JOIN payment p     ON p.resident_id = r.id
WHERE r.deleted_at IS NULL AND p.status = 'verified'
AND EXISTS (
    SELECT 1 FROM document d
    WHERE d.resident_id = r.id AND d.status = 'verified'
) AND similarity(CONCAT_WS(' ', r.first_name, r.second_name, r.last_name), sqlc.arg('query')) > 0.2;

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


-- name: ListUnverifiedResidents :many
SELECT
    sqlc.embed(r), sqlc.embed(p), b.face_url,
    CONCAT_WS(' ', r.first_name, r.second_name, r.last_name) AS full_name,
    CASE
        WHEN EXISTS (
            SELECT 1 FROM document d
            WHERE d.resident_id = r.id AND d.status = 'pending'
        ) THEN 'pending'
        WHEN EXISTS (
            SELECT 1 FROM document d
            WHERE d.resident_id = r.id AND d.status = 'rejected'
        ) THEN 'rejected'
        ELSE 'no documents'
    END AS status
FROM resident r
JOIN biometric b   ON r.id = b.resident_id
JOIN payment p     ON p.resident_id = r.id
WHERE r.deleted_at IS NULL AND p.status = 'verified'
AND NOT EXISTS (
    SELECT 1 FROM document d
    WHERE d.resident_id = r.id AND d.status = 'verified'
)
ORDER BY r.created_at ASC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountListUnverifiedResidents :one
SELECT COUNT(*)
FROM resident r
JOIN payment p     ON p.resident_id = r.id
WHERE r.deleted_at IS NULL AND p.status = 'verified'
AND NOT EXISTS (
    SELECT 1 FROM document d
    WHERE d.resident_id = r.id AND d.status = 'verified'
);

-- name: SearchUnverifiedResidents :many
SELECT
    sqlc.embed(r), sqlc.embed(p), b.face_url,
    CONCAT_WS(' ', r.first_name, r.second_name, r.last_name) AS full_name,
    similarity(CONCAT_WS(' ', r.first_name, r.second_name, r.last_name), sqlc.arg('query')) AS sim,
    CASE
        WHEN EXISTS (
            SELECT 1 FROM document d
            WHERE d.resident_id = r.id AND d.status = 'pending'
        ) THEN 'pending'
        WHEN EXISTS (
            SELECT 1 FROM document d
            WHERE d.resident_id = r.id AND d.status = 'rejected'
        ) THEN 'rejected'
        ELSE 'no documents'
    END AS status
FROM resident r
JOIN biometric b   ON r.id = b.resident_id
JOIN payment p     ON p.resident_id = r.id
WHERE r.deleted_at IS NULL AND p.status = 'verified'
AND NOT EXISTS (
    SELECT 1 FROM document d
    WHERE d.resident_id = r.id AND d.status = 'verified'
) AND similarity(CONCAT_WS(' ', r.first_name, r.second_name, r.last_name), sqlc.arg('query')) > 0.2
ORDER BY r.created_at ASC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountSearchUnverifiedResidents :one
SELECT COUNT(*)
FROM resident r
JOIN payment p     ON p.resident_id = r.id
WHERE r.deleted_at IS NULL AND p.status = 'verified'
AND NOT EXISTS (
    SELECT 1 FROM document d
    WHERE d.resident_id = r.id AND d.status = 'verified'
) AND similarity(CONCAT_WS(' ', r.first_name, r.second_name, r.last_name), sqlc.arg('query')) > 0.2;
