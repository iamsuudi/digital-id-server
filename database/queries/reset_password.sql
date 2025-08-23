-- name: CreatePasswordResetToken :one
INSERT INTO password_reset_tokens (user_id, token, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetValidPasswordResetToken :one
SELECT *
FROM password_reset_tokens
WHERE token = $1
    AND expires_at > NOW()
    AND used_at IS NULL
LIMIT 1;

-- name: MarkTokenAsUsed :exec
UPDATE password_reset_tokens
SET used_at = NOW()
WHERE id = $1;

-- name: UpdateUserPassword :exec
UPDATE account
SET password_hash = $1
WHERE user_id = $2;

-- name: DeleteExpiredTokens :exec
DELETE FROM password_reset_tokens
WHERE expires_at <= NOW()
    OR used_at IS NOT NULL;
