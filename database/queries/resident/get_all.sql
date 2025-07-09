-- name: GetAllResidents :many
SELECT
  id,
  email,
  first_name,
  second_name,
  last_name,
  birth_date,
  gender,
  phone,
  marital_status,
  religion,
  ethnicity,
  disability_status,
  education_level,
  languages_spoken,
  address_id,
  created_at,
  deleted_at
FROM resident
WHERE deleted_at IS NULL
ORDER BY created_at DESC;
