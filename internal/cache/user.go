package cache

import (
	"context"
	"digital-id-server/internal/repository"

	"github.com/google/uuid"
)

type UserReader interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (repository.GetUserByIDRow, error)
}

type UserData map[string]repository.GetUserByIDRow

func (c *Cache) GetUser(ctx context.Context, id uuid.UUID) (repository.GetUserByIDRow, error) {
	key := id.String()

	c.mu.RLock()
	list, ok := c.data.user[key]
	c.mu.RUnlock()
	if ok {
		return list, nil
	}

	list, err := c.dbUser.GetUserByID(ctx, id)
	if err != nil {
		return repository.GetUserByIDRow{}, err
	}

	c.mu.Lock()
	c.data.user[key] = list
	c.mu.Unlock()
	return list, nil
}
