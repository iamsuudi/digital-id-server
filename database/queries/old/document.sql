-- name: CreateDocument :exec
INSERT INTO document (
  type, number, resident_id, url, status
) VALUES ($1, $2, $3, $4, $5);
