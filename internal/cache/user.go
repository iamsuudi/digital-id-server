package cache

import (
	"context"
	"digital-id-server/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
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

func (c *Cache) RequirePlacement() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		raw, _ := ctx.Get("user_id")
		id, _ := raw.(uuid.UUID)

		user, err := c.GetUser(ctx, id)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
			return
		}

		switch user.RoleSlug {
		case "admin":
			{
				if user.CityID != nil {
					ctx.Next()
				} else {
					ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You are not assigned placement"})
					return
				}
			}
		case "manager":
			{
				if user.SubcityID != nil && user.CityID != nil {
					ctx.Next()
				} else {
					ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You are not assigned placement"})
					return
				}
			}
		case "executive":
			{
				if user.KebeleID != nil && user.SubcityID != nil && user.CityID != nil {
					ctx.Next()
				} else {
					ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You are not assigned placement"})
					return
				}
			}
		case "cashier":
			{
				if user.KebeleID != nil && user.SubcityID != nil && user.CityID != nil {
					ctx.Next()
				} else {
					ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You are not assigned placement"})
					return
				}
			}
		case "encoder":
			{
				if user.KebeleID != nil && user.SubcityID != nil && user.CityID != nil {
					ctx.Next()
				} else {
					ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You are not assigned placement"})
					return
				}
			}
		}

		ctx.Next()
	}
}
