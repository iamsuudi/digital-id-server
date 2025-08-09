-- name: CreateBiometric :one
INSERT INTO biometric (resident_id, face_url, blood_type)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateBiometric :exec
UPDATE biometric
SET face_url = $2, blood_type = $3
WHERE resident_id = $1;
