package cache

import (
	"context"
	"digital-id-server/internal/repository"

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
