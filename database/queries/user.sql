-- name: CreateUser :one
INSERT INTO "user" (first_name, second_name, last_name, email, phone, password_hash, role_slug)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *, CONCAT_WS(' ', first_name, second_name, last_name) AS full_name;

-- name: GetUserByEmail :one
SELECT *, CONCAT_WS(' ', first_name, second_name, last_name) AS full_name
FROM "user"
WHERE email = $1 AND deleted_at IS NULL;

-- name: GetUserByID :one
SELECT *, CONCAT_WS(' ', first_name, second_name, last_name) AS full_name
FROM "user"
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetUserScope :one
SELECT city_id, subcity_id, kebele_id FROM "user" WHERE id = $1;

-- name: UpdateUserRole :exec
UPDATE "user"
SET role_slug = $2
WHERE id = $1;

-- name: SoftDeleteUser :exec
UPDATE "user"
SET deleted_at = NOW()
WHERE id = $1;

-- name: ListAllUsers :many
SELECT u.*, r.name AS role,
       CONCAT_WS(' ', first_name, second_name, last_name) AS full_name
FROM "user" u
JOIN role r ON r.slug = u.role_slug
WHERE u.deleted_at IS NULL
ORDER BY u.id
LIMIT  sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CountListUsers :one
SELECT COUNT(*)
FROM "user"
WHERE deleted_at IS NULL;

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
FROM "user"
WHERE deleted_at IS NULL AND
    (sqlc.arg('city_id')::uuid IS NULL OR u.city_id = sqlc.arg('city_id')::uuid) AND
    (sqlc.arg('subcity_id')::uuid IS NULL OR u.subcity_id = sqlc.arg('subcity_id')::uuid) AND
    (sqlc.arg('kebele_id')::uuid IS NULL OR u.kebele_id = sqlc.arg('kebele_id')::uuid);

-- name: SearchUsersUnderScope :many
SELECT *, r.name, CONCAT_WS(' ', first_name, second_name, last_name) AS full_name
FROM "user" u
JOIN role r ON r.slug = u.role_slug
WHERE deleted_at IS NULL AND
    search_vector @@ to_tsquery('english', sqlc.arg('query') || ':*') AND
    (sqlc.arg('city_id')::uuid IS NULL OR u.city_id = sqlc.arg('city_id')::uuid) AND
    (sqlc.arg('subcity_id')::uuid IS NULL OR u.subcity_id = sqlc.arg('subcity_id')::uuid) AND
    (sqlc.arg('kebele_id')::uuid IS NULL OR u.kebele_id = sqlc.arg('kebele_id')::uuid)
ORDER BY
    ts_rank(search_vector, to_tsquery('english', sqlc.arg('query') || ':*')) DESC,
    similarity(first_name || ' ' || second_name || ' ' || last_name, sqlc.arg('query')) DESC,
    created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountUsersSearchUnderScope :one
SELECT COUNT(*)
FROM "user"
WHERE deleted_at IS NULL AND
    search_vector @@ plainto_tsquery('english', sqlc.arg('query')) AND 
    (sqlc.arg('city_id')::uuid IS NULL OR u.city_id = sqlc.arg('city_id')::uuid) AND
    (sqlc.arg('subcity_id')::uuid IS NULL OR u.subcity_id = sqlc.arg('subcity_id')::uuid) AND
    (sqlc.arg('kebele_id')::uuid IS NULL OR u.kebele_id = sqlc.arg('kebele_id')::uuid);
