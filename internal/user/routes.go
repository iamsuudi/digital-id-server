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

	usersGroup := rg.Group("/users", auth.Authenticate())
	{
		usersGroup.GET("/", auth.RolePrefixRedirect("superadmin"), handler.GetAll)
		usersGroup.GET("/role", handler.GetByRole)
		usersGroup.GET("/:id", handler.GetUser)
	}

	superadmin := rg.Group("/superadmin", auth.Authenticate(), auth.Authorize("superadmin"))
	{
		superadmin.GET("/users", handler.RequirePermission("can_edit_price"), handler.GetAllForSuperadmin)
	}
}
