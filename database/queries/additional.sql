-- name: CreateAdditional :exec
INSERT INTO additional (
    resident_id, marital_status, religion, ethnicity, disability, education_level, languages_spoken
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
);