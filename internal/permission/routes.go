package permission

import (
	"digital-id-server/internal/auth"
	"digital-id-server/internal/cache"
	"digital-id-server/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(rg *gin.RouterGroup, db *pgxpool.Pool, q *repository.Queries, cache *cache.Cache) {
	service := NewService(db, q)
	handler := NewHandler(service, cache)

	r := rg.Group("/permissions", auth.Authenticate())
	{
		r.GET("/", handler.GetAllPermissions)
		r.POST("/", handler.CreatePermission)
		r.GET("/assignable", handler.GetAssignablePermissions)
		r.GET("/user", handler.GetUniversalPermissionsForUser)
		r.PUT("/override/set", handler.OverridePermission)
		r.DELETE("/override/remove", handler.RemoveOverride)
		r.DELETE("/:name", handler.DeletePermission)
	}
}
