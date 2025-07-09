-- name: CreateEmployment :exec
INSERT INTO employment (
  resident_id, status, occupation, employer_name, work_address
) VALUES ($1, $2, $3, $4, $5);
