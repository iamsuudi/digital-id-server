package cache

import (
	"context"
	"digital-id-server/internal/repository"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PermsReader interface {
	GetEffectivePermissions(ctx context.Context, id uuid.UUID) ([]repository.GetEffectivePermissionsForUserRow, error)
}

type PermsData map[string][]repository.GetEffectivePermissionsForUserRow

func (c *Cache) GetPerms(ctx context.Context, id uuid.UUID) ([]repository.GetEffectivePermissionsForUserRow, error) {
	key := id.String()

	c.mu.RLock()
	list, ok := c.data.perms[key]
	c.mu.RUnlock()
	if ok {
		return list, nil
	}

	list, err := c.dbPerms.GetEffectivePermissions(ctx, id)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	c.data.perms[key] = list
	c.mu.Unlock()
	return list, nil
}

func(c *Cache) RequirePerms(required ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		raw, _ := ctx.Get("user_id")
		id, _ := raw.(uuid.UUID)
		
		perms, err := c.GetPerms(ctx, id)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Permissions unavailable"})
			return
		}

		for _, reqPerm := range required {
			if !slices.ContainsFunc(perms, func(p repository.GetEffectivePermissionsForUserRow) bool {
				return p.Name == reqPerm
			}) {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
				return
			}
		}
		ctx.Next()
	}
}
