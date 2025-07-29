-- name: CreateUser :one
INSERT INTO "user" (first_name, second_name, last_name, email, phone, password_hash, role_slug)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *, CONCAT_WS(' ', first_name, second_name, last_name) AS full_name;

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

-- name: UpdateUserRole :exec
UPDATE "user"
SET role_slug = $2
WHERE id = $1;

-- name: UpdateUserPlacement :exec
WITH clear_old AS (
    UPDATE "user"
    SET city_id = NULL, subcity_id = NULL, kebele_id = NULL
    WHERE city_id = $2
)
UPDATE "user" AS U
SET city_id = $2, subcity_id = $3, kebele_id = $4
WHERE u.id = $1;

-- name: SoftDeleteUser :exec
UPDATE "user"
SET deleted_at = NOW()
WHERE id = $1;

-- name: ListAllUsers :many
SELECT u.*, CONCAT_WS(' ', u.first_name, u.second_name, u.last_name) AS full_name,
    c.name AS city_name, sc.name AS subcity_name, k.name AS kebele_name,
    r.name AS role_name, r.level_rank AS role_level_rank
FROM "user" u
LEFT JOIN city c ON c.id = u.city_id
LEFT JOIN subcity sc ON sc.id = u.subcity_id
LEFT JOIN kebele k ON k.id = u.kebele_id
LEFT JOIN role r ON r.slug = u.role_slug
WHERE u.deleted_at IS NULL
ORDER BY u.created_at ASC
LIMIT  sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountListUsers :one
SELECT COUNT(*)
FROM "user"
WHERE deleted_at IS NULL;

-- name: SearchUsers :many
SELECT u.*, CONCAT_WS(' ', u.first_name, u.second_name, u.last_name) AS full_name,
    c.name AS city_name, sc.name AS subcity_name, k.name AS kebele_name,
    r.name AS role_name, r.level_rank AS role_level_rank,
    similarity(CONCAT_WS(' ', first_name, second_name, last_name), sqlc.arg('query')) AS sim
FROM "user" u
LEFT JOIN city c ON c.id = u.city_id
LEFT JOIN subcity sc ON sc.id = u.subcity_id
LEFT JOIN kebele k ON k.id = u.kebele_id
LEFT JOIN role r ON r.slug = u.role_slug
WHERE u.deleted_at IS NULL AND
    similarity(CONCAT_WS(' ', first_name, second_name, last_name), sqlc.arg('query')) > 0.2
ORDER BY sim DESC, u.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountUsersSearch :one
SELECT COUNT(*)
FROM "user"
WHERE deleted_at IS NULL AND
    similarity(CONCAT_WS(' ', first_name, second_name, last_name), sqlc.arg('query')) > 0.2;

-- name: ListUsersUnderScope :many
SELECT u.*, r.name AS role, CONCAT_WS(' ', first_name, second_name, last_name) AS full_name
FROM "user" u
JOIN role r ON r.slug = u.role_slug
WHERE deleted_at IS NULL AND
    (sqlc.arg('city_id')::uuid IS NULL OR u.city_id = sqlc.arg('city_id')::uuid) AND
    (sqlc.arg('subcity_id')::uuid IS NULL OR u.subcity_id = sqlc.arg('subcity_id')::uuid) AND
    (sqlc.arg('kebele_id')::uuid IS NULL OR u.kebele_id = sqlc.arg('kebele_id')::uuid)
ORDER BY u.id
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountListUsersUnderScope :one
SELECT COUNT(*)
FROM "user" u
WHERE deleted_at IS NULL AND
    (sqlc.arg('city_id')::uuid IS NULL OR u.city_id = sqlc.arg('city_id')::uuid) AND
    (sqlc.arg('subcity_id')::uuid IS NULL OR u.subcity_id = sqlc.arg('subcity_id')::uuid) AND
    (sqlc.arg('kebele_id')::uuid IS NULL OR u.kebele_id = sqlc.arg('kebele_id')::uuid);

-- name: SearchUsersUnderScope :many
SELECT *, r.name, CONCAT_WS(' ', first_name, second_name, last_name) AS full_name,
    similarity(CONCAT_WS(' ', first_name, second_name, last_name), sqlc.arg('query')) AS sim
FROM "user" u
JOIN role r ON r.slug = u.role_slug
WHERE deleted_at IS NULL AND
    similarity(CONCAT_WS(' ', first_name, second_name, last_name), sqlc.arg('query')) > 0.2 AND
    (sqlc.arg('city_id')::uuid IS NULL OR u.city_id = sqlc.arg('city_id')::uuid) AND
    (sqlc.arg('subcity_id')::uuid IS NULL OR u.subcity_id = sqlc.arg('subcity_id')::uuid) AND
    (sqlc.arg('kebele_id')::uuid IS NULL OR u.kebele_id = sqlc.arg('kebele_id')::uuid)
ORDER BY sim DESC, created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountUsersSearchUnderScope :one
SELECT COUNT(*)
FROM "user" u
WHERE deleted_at IS NULL AND
    similarity(CONCAT_WS(' ', first_name, second_name, last_name), sqlc.arg('query')) > 0.2 AND
    (sqlc.arg('city_id')::uuid IS NULL OR u.city_id = sqlc.arg('city_id')::uuid) AND
    (sqlc.arg('subcity_id')::uuid IS NULL OR u.subcity_id = sqlc.arg('subcity_id')::uuid) AND
    (sqlc.arg('kebele_id')::uuid IS NULL OR u.kebele_id = sqlc.arg('kebele_id')::uuid);


-- name: ListByRole :many
SELECT *, CONCAT_WS(' ', first_name, second_name, last_name) AS full_name
FROM "user"
WHERE role_slug = sqlc.arg('role_slug') AND deleted_at IS NULL
ORDER BY created_at ASC
LIMIT  sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountListByRole :one
SELECT COUNT(*)
FROM "user"
WHERE role_slug = sqlc.arg('role_slug') AND deleted_at IS NULL;

-- name: SearchByRole :many
SELECT *, CONCAT_WS(' ', first_name, second_name, last_name) AS full_name,
    similarity(CONCAT_WS(' ', first_name, second_name, last_name), sqlc.arg('query')) AS sim
FROM "user"
WHERE role_slug = sqlc.arg('role_slug') AND deleted_at IS NULL AND
    similarity(CONCAT_WS(' ', first_name, second_name, last_name), sqlc.arg('query')) > 0.2
ORDER BY sim DESC, created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountByRoleSearch :one
SELECT COUNT(*)
FROM "user"
WHERE role_slug = sqlc.arg('role_slug') AND deleted_at IS NULL AND
    similarity(CONCAT_WS(' ', first_name, second_name, last_name), sqlc.arg('query')) > 0.2;
