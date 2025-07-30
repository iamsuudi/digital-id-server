package cache

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type Data struct {
	user  UserData
	perms PermsData
}

type Cache struct {
	dbUser  UserReader
	dbPerms PermsReader
	mu      sync.RWMutex
	data    Data
	ttl     time.Duration
}

func New(dbU UserReader, dbP PermsReader) *Cache {
	return &Cache{
		dbUser:  dbU,
		dbPerms: dbP,
		data: Data{
			user:  make(UserData),
			perms: make(PermsData),
		},
		ttl: 3 * time.Minute,
	}
}

func (c *Cache) Invalidate(id uuid.UUID) {
	c.mu.Lock()
	delete(c.data.user, id.String())
	delete(c.data.perms, id.String())
	c.mu.Unlock()
}
