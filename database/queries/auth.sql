-- name: CanActorTouchTarget :one
SELECT CASE
    WHEN $1 = $2 THEN false
    WHEN a.role_slug = t.role_slug THEN false
    WHEN (SELECT level_rank FROM role WHERE slug = a.role_slug)
         >= (SELECT level_rank FROM role WHERE slug = t.role_slug) THEN false
    WHEN (a.city_id     IS NOT NULL AND a.city_id     = t.city_id) OR
         (a.subcity_id IS NOT NULL AND a.subcity_id = t.subcity_id) OR
         (a.kebele_id IS NOT NULL AND a.kebele_id = t.kebele_id) THEN true
    ELSE false
END AS ok
FROM "user" a, "user" t
WHERE a.id = $1 AND t.id = $2;

-- name: CanActorManipulateRole :one
-- returns true if the actor’s role is strictly more powerful than target_role
SELECT CASE
    WHEN actor.level_rank = target.level_rank THEN false   -- same role
    WHEN actor.level_rank < target.level_rank THEN true
    ELSE false
END AS ok
FROM role AS actor, role AS target
WHERE actor.slug = (SELECT role_slug FROM "user" WHERE id = sqlc.arg('id'))
  AND target.slug = sqlc.arg('role');

-- name: CanActorGrantPermissionToRole :one
SELECT CASE
      WHEN (
          -- 1. actor role is strictly higher than target role
          (SELECT level_rank FROM role WHERE slug = (
              SELECT role_slug FROM "user" u WHERE u.id = sqlc.arg('id')
          )) <
          (SELECT level_rank FROM role r WHERE r.slug = sqlc.arg('role'))
      ) AND (
          -- 2. permission exists in role(s) at or below actor’s level
          EXISTS (
              SELECT 1
              FROM role_permission rp
              JOIN role r ON r.slug = rp.role_slug
              WHERE rp.permission_name = sqlc.arg('permission')
                AND r.level_rank >= (
                    SELECT level_rank FROM role WHERE slug = (
                        SELECT role_slug FROM "user" WHERE id = sqlc.arg('id')
                    )
                )
          )
      )
      THEN true
      ELSE false
  END AS ok;
