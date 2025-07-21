package auth

import (
	"net/http"
	"slices"
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

func hasRole(ctx *gin.Context, expectedRole string) bool {
	raw, exists := ctx.Get(ContextUserRoleKey)
	role, ok := raw.(string)
	if !exists || !ok {
		return false
	}
	return role == expectedRole
}

func Authorize(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		raw, exists := ctx.Get(ContextUserRoleKey)
		role, ok := raw.(string)
		if !exists || !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		if slices.Contains(allowedRoles, role) {
			ctx.Next()
			return
		}

		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
	}
}

// RolePrefixRedirect returns a Gin middleware that internally routes
//
//	/foo/bar   ->   /{prefix}/foo/bar
//
// when the current user has the required role.
func RolePrefixRedirect(prefix string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if hasRole(ctx, prefix) {
			// Construct the new path with prefix
			newPath := "/api/v1/" + prefix + strings.TrimPrefix(ctx.Request.URL.Path, "/api/v1") + "?" + ctx.Request.URL.RawQuery
			ctx.Redirect(http.StatusFound, newPath) // 302 redirect
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
