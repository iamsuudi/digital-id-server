-- name: CreateAccount :exec
INSERT INTO account (user_id, password_hash)
VALUES ($1, $2);

-- name: GetAccount :one
SELECT * FROM account
WHERE account.user_id = $1;

-- name: CreateUser :one
INSERT INTO "user" (first_name, second_name, last_name, email, phone, role_slug)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *, CONCAT_WS(' ', first_name, second_name, last_name) AS full_name;

-- name: UpdateUserInfo :one
UPDATE "user"
SET first_name = $2, second_name = $3, last_name = $4, email = $5, phone = $6
WHERE id = $1
RETURNING *;

-- name: GetUserByEmail :one
SELECT u.*, CONCAT_WS(' ', u.first_name, u.second_name, u.last_name) AS full_name,
    c.name AS city_name, sc.name AS subcity_name, k.name AS kebele_name,
    r.name AS role_name, r.level_rank AS role_level_rank
FROM "user" u
LEFT JOIN role r ON r.slug = u.role_slug
LEFT JOIN city c ON c.id = u.city_id
LEFT JOIN subcity sc ON sc.id = u.subcity_id
LEFT JOIN kebele k ON k.id = u.kebele_id
WHERE u.email = $1 AND u.deleted_at IS NULL;

-- name: GetUserByID :one
SELECT u.*, CONCAT_WS(' ', u.first_name, u.second_name, u.last_name) AS full_name,
    c.name AS city_name, sc.name AS subcity_name, k.name AS kebele_name,
    r.name AS role_name, r.level_rank AS role_level_rank
FROM "user" u
LEFT JOIN role r ON r.slug = u.role_slug
LEFT JOIN city c ON c.id = u.city_id
LEFT JOIN subcity sc ON sc.id = u.subcity_id
LEFT JOIN kebele k ON k.id = u.kebele_id
WHERE u.id = $1 AND u.deleted_at IS NULL;

-- name: GetUserScope :one
SELECT city_id, subcity_id, kebele_id FROM "user" WHERE id = $1;

-- name: GetUserRole :one
SELECT role_slug FROM "user" WHERE id = $1;

-- name: UpdateUserRole :exec
UPDATE "user"
SET role_slug = $2
WHERE id = $1;

-- name: GrantUserPlacement :exec
UPDATE "user"
SET city_id = sqlc.narg('city_id'),
    subcity_id = sqlc.narg('subcity_id'),
    kebele_id = sqlc.narg('kebele_id')
WHERE id = sqlc.arg('id');

-- name: RevokeUserPlacement :exec
UPDATE "user"
SET city_id = sqlc.narg('city_id')::uuid,
    subcity_id = sqlc.narg('subcity_id')::uuid,
    kebele_id = sqlc.narg('kebele_id')::uuid
WHERE id = sqlc.arg('id');

-- name: SoftDeleteUser :exec
UPDATE "user"
SET deleted_at = NOW()
WHERE id = $1;

-- name: ListUsersByKebeleAndRole :many
SELECT id, CONCAT_WS(' ', first_name, second_name, last_name) AS full_name
FROM "user"
WHERE kebele_id = $1
  AND role_slug = $2
  AND deleted_at IS NULL;

-- name: ListUsersUnderScope :many
SELECT u.*, r.name AS role_name, CONCAT_WS(' ', u.first_name, u.second_name, u.last_name) AS full_name,
    c.name AS city_name, sc.name AS subcity_name, k.name AS kebele_name,
    r.name AS role_name, r.level_rank AS role_level_rank
FROM "user" u
JOIN role r ON r.slug = u.role_slug
LEFT JOIN city c ON c.id = u.city_id
LEFT JOIN subcity sc ON sc.id = u.subcity_id
LEFT JOIN kebele k ON k.id = u.kebele_id
WHERE u.deleted_at IS NULL AND
    sqlc.arg('rank') < r.level_rank AND
    (sqlc.narg('city_id')::uuid IS NULL OR u.city_id = sqlc.narg('city_id')::uuid) AND
    (sqlc.narg('subcity_id')::uuid IS NULL OR u.subcity_id = sqlc.narg('subcity_id')::uuid) AND
    (sqlc.narg('kebele_id')::uuid IS NULL OR u.kebele_id = sqlc.narg('kebele_id')::uuid)
ORDER BY u.created_at ASC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountListUsersUnderScope :one
SELECT COUNT(*)
FROM "user" u
JOIN role r ON  r.slug = u.role_slug
WHERE deleted_at IS NULL AND
    sqlc.arg('rank') < r.level_rank AND
    (sqlc.narg('city_id')::uuid IS NULL OR u.city_id = sqlc.narg('city_id')::uuid) AND
    (sqlc.narg('subcity_id')::uuid IS NULL OR u.subcity_id = sqlc.narg('subcity_id')::uuid) AND
    (sqlc.narg('kebele_id')::uuid IS NULL OR u.kebele_id = sqlc.narg('kebele_id')::uuid);

-- name: SearchUsersUnderScope :many
SELECT *, r.name, CONCAT_WS(' ', u.first_name, u.second_name, u.last_name) AS full_name,
    c.name AS city_name, sc.name AS subcity_name, k.name AS kebele_name,
    r.name AS role_name, r.level_rank AS role_level_rank,
    similarity(CONCAT_WS(' ', u.first_name, u.second_name, u.last_name), sqlc.arg('query')) AS sim
FROM "user" u
JOIN role r ON r.slug = u.role_slug
LEFT JOIN city c ON c.id = u.city_id
LEFT JOIN subcity sc ON sc.id = u.subcity_id
LEFT JOIN kebele k ON k.id = u.kebele_id
WHERE u.deleted_at IS NULL AND
    sqlc.arg('rank') < r.level_rank AND
    similarity(CONCAT_WS(' ', u.first_name, u.second_name, u.last_name), sqlc.arg('query')) > 0.2 AND
    (sqlc.narg('city_id')::uuid IS NULL OR u.city_id = sqlc.narg('city_id')::uuid) AND
    (sqlc.narg('subcity_id')::uuid IS NULL OR u.subcity_id = sqlc.narg('subcity_id')::uuid) AND
    (sqlc.narg('kebele_id')::uuid IS NULL OR u.kebele_id = sqlc.narg('kebele_id')::uuid)
ORDER BY sim DESC, u.created_at ASC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountSearchUsersUnderScope :one
SELECT COUNT(*)
FROM "user" u
JOIN role r ON r.slug = u.role_slug
WHERE deleted_at IS NULL AND
    sqlc.arg('rank') < r.level_rank AND
    similarity(CONCAT_WS(' ', u.first_name, u.second_name, u.last_name), sqlc.arg('query')) > 0.2 AND
    (sqlc.narg('city_id')::uuid IS NULL OR u.city_id = sqlc.narg('city_id')::uuid) AND
    (sqlc.narg('subcity_id')::uuid IS NULL OR u.subcity_id = sqlc.narg('subcity_id')::uuid) AND
    (sqlc.narg('kebele_id')::uuid IS NULL OR u.kebele_id = sqlc.narg('kebele_id')::uuid);

-- name: ListUsersByRole :many
SELECT *, CONCAT_WS(' ', first_name, second_name, last_name) AS full_name
FROM "user"
WHERE role_slug = sqlc.arg('role_slug') AND deleted_at IS NULL
ORDER BY created_at ASC
LIMIT  sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountListUsersByRole :one
SELECT COUNT(*)
FROM "user"
WHERE role_slug = sqlc.arg('role_slug') AND deleted_at IS NULL;

-- name: SearchUsersByRole :many
SELECT *, CONCAT_WS(' ', first_name, second_name, last_name) AS full_name,
    similarity(CONCAT_WS(' ', first_name, second_name, last_name), sqlc.arg('query')) AS sim
FROM "user"
WHERE role_slug = sqlc.arg('role_slug') AND deleted_at IS NULL AND
    similarity(CONCAT_WS(' ', first_name, second_name, last_name), sqlc.arg('query')) > 0.2
ORDER BY sim DESC, created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountSearchUsersByRole :one
SELECT COUNT(*)
FROM "user"
WHERE role_slug = sqlc.arg('role_slug') AND deleted_at IS NULL AND
    similarity(CONCAT_WS(' ', first_name, second_name, last_name), sqlc.arg('query')) > 0.2;
