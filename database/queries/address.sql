-- name: GetAddress :one
SELECT * FROM address
WHERE house_number = $1 AND kebele_id = $2 AND city_id = $3
LIMIT 1;

-- name: CreateAddress :one
INSERT INTO address (
  house_number, kebele_id, city_id
) VALUES ($1, $2, $3)
RETURNING id;
