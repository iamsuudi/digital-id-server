-- name: CreateResident :one
INSERT INTO resident (
  email, first_name, second_name, last_name, birth_date, gender, phone,
  marital_status, religion, ethnicity, disability_status, education_level,
  languages_spoken, address_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7,
  $8, $9, $10, $11, $12, $13, $14
)
RETURNING id;

-- name: UpdateResidentAddress :exec
UPDATE resident SET address_id = $1 WHERE id = $2;
 