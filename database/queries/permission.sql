-- name: ListPermissions :many
SELECT name, label, description FROM permission ORDER BY name;

-- name: GrantPermissionToRole :exec
INSERT INTO role_permission (role_slug, permission_name)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: RevokePermissionFromRole :exec
DELETE FROM role_permission
WHERE role_slug = $1 AND permission_name = $2;

-- name: GetEffectivePermissionsForUser :many
WITH RECURSIVE role_tree AS (
    -- start from the user’s role
    SELECT r.slug, r.level_rank
    FROM role r
    WHERE r.slug = (SELECT role_slug FROM "user" WHERE id = $1)
    UNION
    SELECT r.slug, r.level_rank
    FROM role r
    JOIN role_tree rt ON r.parent_role_slug = rt.slug
),
inherited AS (
    SELECT DISTINCT rp.permission_name
    FROM role_permission rp
    WHERE rp.role_slug IN (SELECT slug FROM role_tree)
),
overrides AS (
    SELECT permission_name, is_granted
    FROM user_permission_override
    WHERE user_id = $1
)
SELECT
    p.name,
    p.label,
    p.description,
    COALESCE(o.is_granted, true) AS effective,
    o.is_granted IS NOT NULL     AS overridden
FROM permission p
LEFT JOIN inherited i ON i.permission_name = p.name
LEFT JOIN overrides o ON o.permission_name = p.name
WHERE o.is_granted = true
   OR (o.is_granted IS NULL AND i.permission_name IS NOT NULL)
ORDER BY p.name;

-- name: GetPermissionOverridesForUser :many
SELECT permission_name, is_granted
FROM user_permission_override
WHERE user_id = $1;

-- name: SetUserPermissionOverride :exec
INSERT INTO user_permission_override (user_id, permission_name, is_granted, granted_by)
VALUES ($1, $2, $3, $4)
ON CONFLICT (user_id, permission_name)
DO UPDATE SET is_granted = EXCLUDED.is_granted, granted_by = EXCLUDED.granted_by, granted_at = now();

-- name: RemoveUserPermissionOverride :exec
DELETE FROM user_permission_override
WHERE user_id = $1 AND permission_name = $2;

-- name: ListPermissionOverridesForUser :many
SELECT permission_name, is_granted, granted_by, granted_at
FROM user_permission_override
WHERE user_id = $1
ORDER BY permission_name;

-- name: GetUniversalPermissionMatrixForUser :many
SELECT
    p.name,
    p.label,
    p.description,

    -- 1. does the TARGET ROLE itself have the permission?
    EXISTS (
        SELECT 1
        FROM role_permission trp
        WHERE trp.role_slug = sqlc.arg('target_role_slug')
          AND trp.permission_name = p.name
    ) AS effective,

    -- 2. user-level override for this permission
    COALESCE(
        (SELECT o.is_granted
         FROM user_permission_override o
         WHERE o.user_id = sqlc.arg('target_user_id')
           AND o.permission_name = p.name),
        NULL
    ) AS overridden,

    -- 3. does the ACTOR ROLE itself have the permission?
    EXISTS (
        SELECT 1
        FROM role_permission arp
        WHERE arp.role_slug = sqlc.arg('actor_role_slug')
          AND arp.permission_name = p.name
    ) AS grantable

FROM permission p
ORDER BY p.name;
