package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserIDKey   = "user_id"
	ContextUserRoleKey = "user_role"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("oict_jwt")
		if err != nil || strings.TrimSpace(token) == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: missing token"})
			return
		}

		claims, err := ParseJWT(token)
		if err != nil || claims == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: invalid or expired token"})
			return
		}

		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextUserRoleKey, claims.Role)
		c.Next()
	}
}

func Authorize(allowedRoles ...string) gin.HandlerFunc {
	roleSet := make(map[string]bool, len(allowedRoles))
	for _, role := range allowedRoles {
		roleSet[role] = true
	}
	return func(ctx *gin.Context) {
		raw, exists := ctx.Get(ContextUserRoleKey)
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Access Denied: Role missing",
			})
		}

		role, ok := raw.(string)
		if !ok || !roleSet[role] {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Access Denied: Not Authorized",
			})
		}
	}
}
