package user

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

	r := rg.Group("/users", auth.Authenticate())
	{
		r.GET("/", cache.RequirePlacement(), handler.GetUsers)
		r.GET("/role", handler.GetUsersByRole)
		r.GET("/:id", handler.GetUser)
		r.PUT("/:id/role", handler.UpdateUserRole)
		r.PUT("/:id/info", handler.UpdateUserInfo)
	}
}
