-- name: CreateBiometric :one
INSERT INTO biometric (resident_id, face_url, blood_type)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateBiometric :exec
UPDATE biometric
SET
    face_url = COALESCE(sqlc.narg('face_url'), face_url),
    blood_type = COALESCE(sqlc.narg('blood_type'), blood_type)
WHERE resident_id = $1;
