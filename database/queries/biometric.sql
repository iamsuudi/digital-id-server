-- name: CreateBiometric :one
INSERT INTO biometric (resident_id, face_url, blood_type)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetBiometric :one
SELECT * FROM biometric WHERE resident_id = $1;

-- name: UpdateBiometric :one
UPDATE biometric
SET
    face_url = COALESCE(sqlc.narg('face_url'), face_url),
    blood_type = COALESCE(sqlc.narg('blood_type'), blood_type)
WHERE resident_id = $1
RETURNING *;
