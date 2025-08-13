-- name: GetSetting :one
SELECT *
FROM setting
WHERE id = $1;

-- name: UpdateSetting :exec
UPDATE setting
SET idcard_expiration_duration = $1
WHERE id = $1;
