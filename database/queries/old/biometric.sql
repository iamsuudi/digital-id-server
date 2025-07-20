-- name: CreateBiometric :exec
INSERT INTO biometric (
  resident_id, fingerprint, blood_type, face
) VALUES ($1, $2, $3, $4);
