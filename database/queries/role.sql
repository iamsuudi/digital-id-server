-- name: ListRoles :many
SELECT * FROM role ORDER BY level_rank;

-- name: GetAssignableRolesForActor :many
SELECT slug, name, level_rank
FROM role
WHERE level_rank > (
        SELECT r.level_rank FROM role r WHERE r.slug = $1
)
ORDER BY level_rank;

-- name: GetCurrentUserMaxRoleLevel :one
SELECT r.level_rank
FROM "user" u
JOIN role r ON r.slug = u.role_slug
WHERE u.id = $1;

-- name: GetAssignablePermissionsForActor :many
SELECT p.name, p.label, p.description
FROM permission p
WHERE EXISTS (
    SELECT 1
    FROM role_permission rp
    JOIN role role ON role.slug = rp.role_slug
    WHERE rp.permission_name = p.name
      AND role.level_rank = (
          SELECT level_rank FROM role WHERE slug = (SELECT role_slug FROM "user" WHERE id = $1)
      )
)
ORDER BY p.name;

-- name: ListPermissionMatrixInScope :many
WITH scope_users AS (
    SELECT u.id, u.email, r.slug AS role_slug, r.name AS role_name
    FROM "user" u
    JOIN role r ON r.slug = u.role_slug
    WHERE
    (sqlc.arg('city_id')::uuid IS NULL OR u.city_id = sqlc.arg('city_id')::uuid) AND
    (sqlc.arg('subcity_id')::uuid IS NULL OR u.subcity_id = sqlc.arg('subcity_id')::uuid) AND
    (sqlc.arg('kebele_id')::uuid IS NULL OR u.kebele_id = sqlc.arg('kebele_id')::uuid)
)
SELECT
    su.id             AS user_id,
    su.email,
    su.role_slug,
    su.role_name,
    p.name            AS permission_name,
    p.label,
    p.description,
    COALESCE(
        o.is_granted,
        EXISTS (
            SELECT 1
            FROM role_permission rp
            JOIN role child ON child.slug = rp.role_slug
            WHERE rp.permission_name = p.name
              AND child.level_rank >= (
                  SELECT level_rank FROM role WHERE slug = su.role_slug
              )
        )
    ) AS granted,
    o.is_granted IS NOT NULL AS overridden
FROM permission p
CROSS JOIN scope_users su
LEFT  JOIN user_permission_override o
       ON o.user_id = su.id AND o.permission_name = p.name
ORDER BY su.role_slug, su.id, p.name;

-- name: ListRolePermissions :many
SELECT
    r.slug          AS role_slug,
    r.name          AS role_name,
    r.level_rank    AS role_level_rank,
    p.name          AS permission_name,
    p.label         AS permission_label,
    p.description   AS permission_description,
    CASE
        WHEN rp.permission_name IS NOT NULL THEN true
        ELSE false
    END             AS assigned
FROM role r
CROSS JOIN permission p                   -- every permission
LEFT JOIN role_permission rp
       ON rp.role_slug = r.slug
      AND rp.permission_name = p.name
ORDER BY r.level_rank, p.name;
