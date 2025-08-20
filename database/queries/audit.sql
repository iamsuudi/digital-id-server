-- name: InsertAuditLog :exec
INSERT INTO audit_log (
    actor_id, target_user_id, target_role_slug, target_resident_id, target_kebele_id,
    target_subcity_id, target_city_id, action_type, object_type, object_id, diff
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9,$10, $11);

-- name: GetAuditLog :one
-- Return audit log under scope
SELECT log.*, ar.name AS actor_role,
    CONCAT_WS(' ', au.first_name, au.second_name, au.last_name) AS actor_name,
    CONCAT_WS(' ', tu.first_name, tu.second_name, tu.last_name) AS target_user_name
FROM audit_log log
JOIN "user" au ON au.id = log.actor_id
JOIN role ar ON ar.slug = au.role_slug
LEFT JOIN city c ON c.id = au.city_id
LEFT JOIN subcity sc ON sc.id = au.subcity_id
LEFT JOIN kebele k ON k.id = au.kebele_id
LEFT JOIN "user" tu ON tu.id = log.target_user_id
LEFT JOIN role r ON r.slug = log.target_role_slug
WHERE log.actor_id IS NOT NULL AND log.id = sqlc.arg('id') AND
    (sqlc.narg('city_id')::uuid IS NULL OR au.city_id = sqlc.narg('city_id')::uuid) AND
    (sqlc.narg('subcity_id')::uuid IS NULL OR au.subcity_id = sqlc.narg('subcity_id')::uuid) AND
    (sqlc.narg('kebele_id')::uuid IS NULL OR au.kebele_id = sqlc.narg('kebele_id')::uuid);

-- name: ListAuditLogs :many
-- Returns audit logs under scope
SELECT log.*, ar.name AS actor_role,
    CONCAT_WS(' ', au.first_name, au.second_name, au.last_name) AS actor_name,
    CONCAT_WS(' ', tu.first_name, tu.second_name, tu.last_name) AS target_user_name
FROM audit_log log
JOIN "user" au ON au.id = log.actor_id
JOIN role ar ON ar.slug = au.role_slug
LEFT JOIN city c ON c.id = au.city_id
LEFT JOIN subcity sc ON sc.id = au.subcity_id
LEFT JOIN kebele k ON k.id = au.kebele_id
LEFT JOIN "user" tu ON tu.id = log.target_user_id
LEFT JOIN role r ON r.slug = log.target_role_slug
WHERE log.actor_id IS NOT NULL AND
    (sqlc.narg('city_id')::uuid IS NULL OR au.city_id = sqlc.narg('city_id')::uuid) AND
    (sqlc.narg('subcity_id')::uuid IS NULL OR au.subcity_id = sqlc.narg('subcity_id')::uuid) AND
    (sqlc.narg('kebele_id')::uuid IS NULL OR au.kebele_id = sqlc.narg('kebele_id')::uuid)
ORDER BY log.ts DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountListAuditLogs :one
SELECT COUNT(*)
FROM audit_log log
JOIN "user" u ON u.id = log.actor_id
WHERE log.actor_id IS NOT NULL AND
    (sqlc.narg('city_id')::uuid IS NULL OR u.city_id = sqlc.narg('city_id')::uuid) AND
    (sqlc.narg('subcity_id')::uuid IS NULL OR u.subcity_id = sqlc.narg('subcity_id')::uuid) AND
    (sqlc.narg('kebele_id')::uuid IS NULL OR u.kebele_id = sqlc.narg('kebele_id')::uuid);

-- name: ListUserAuditLogs :many
SELECT log.*, ar.name AS actor_role, CONCAT_WS(' ', au.first_name, au.second_name, au.last_name) AS actor_name
FROM audit_log log
JOIN "user" tu ON tu.id = log.target_user_id
JOIN "user" au ON au.id = log.actor_id
JOIN role ar ON ar.slug = au.role_slug
WHERE log.target_user_id = sqlc.arg('id')
ORDER BY log.ts DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountListUserAuditLogs :one
SELECT COUNT(*)
FROM audit_log log
WHERE log.target_user_id = sqlc.arg('id');

-- name: ListResidentAuditLogs :many
SELECT log.*, ar.name AS actor_role, CONCAT_WS(' ', au.first_name, au.second_name, au.last_name) AS actor_name
FROM audit_log log
JOIN resident tr ON tr.id = log.target_resident_id
JOIN "user" au ON au.id = log.actor_id
JOIN role ar ON ar.slug = au.role_slug
WHERE log.target_resident_id = sqlc.arg('id')
ORDER BY log.ts DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountListResidentAuditLogs :one
SELECT COUNT(*)
FROM audit_log log
WHERE log.target_resident_id = sqlc.arg('id');

-- name: ListCityAuditLogs :many
SELECT log.*, ar.name AS actor_role, CONCAT_WS(' ', au.first_name, au.second_name, au.last_name) AS actor_name
FROM audit_log log
JOIN city t ON t.id = log.target_city_id
JOIN "user" au ON au.id = log.actor_id
JOIN role ar ON ar.slug = au.role_slug
WHERE log.target_city_id = sqlc.arg('id')
ORDER BY log.ts DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountListCityAuditLogs :one
SELECT COUNT(*)
FROM audit_log log
WHERE log.target_city_id = sqlc.arg('id');

-- name: ListSubCityAuditLogs :many
SELECT log.*, ar.name AS actor_role, CONCAT_WS(' ', au.first_name, au.second_name, au.last_name) AS actor_name
FROM audit_log log
JOIN subcity t ON t.id = log.target_subcity_id
JOIN "user" au ON au.id = log.actor_id
JOIN role ar ON ar.slug = au.role_slug
WHERE log.target_subcity_id = sqlc.arg('id')
ORDER BY log.ts DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountListSubCityAuditLogs :one
SELECT COUNT(*)
FROM audit_log log
WHERE log.target_subcity_id = sqlc.arg('id');

-- name: ListKebeleAuditLogs :many
SELECT log.*, ar.name AS actor_role, CONCAT_WS(' ', au.first_name, au.second_name, au.last_name) AS actor_name
FROM audit_log log
JOIN kebele t ON t.id = log.target_kebele_id
JOIN "user" au ON au.id = log.actor_id
JOIN role ar ON ar.slug = au.role_slug
WHERE log.target_kebele_id = sqlc.arg('id')
ORDER BY log.ts DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountListKebeleAuditLogs :one
SELECT COUNT(*)
FROM audit_log log
WHERE log.target_kebele_id = sqlc.arg('id');

-- name: ListRoleAuditLogs :many
SELECT log.*, ar.name AS actor_role, CONCAT_WS(' ', au.first_name, au.second_name, au.last_name) AS actor_name
FROM audit_log log
JOIN role t ON t.slug = log.target_role_slug
JOIN "user" au ON au.id = log.actor_id
JOIN role ar ON ar.slug = au.role_slug
WHERE log.target_role_slug = sqlc.arg('slug')
ORDER BY log.ts DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountListRoleAuditLogs :one
SELECT COUNT(*)
FROM audit_log log
WHERE log.target_role_slug = sqlc.arg('slug');
