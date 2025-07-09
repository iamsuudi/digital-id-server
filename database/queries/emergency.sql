-- name: CreateEmergency :exec
INSERT INTO emergency (
  resident_id, name, relation, phone
) VALUES ($1, $2, $3, $4);
