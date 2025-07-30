package role

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"digital-id-server/internal/auth"
	"digital-id-server/internal/repository"
)

func RegisterRoutes(rg *gin.RouterGroup, db *pgxpool.Pool, q *repository.Queries) {
	service := NewService(db, q)
	handler := NewHandler(service)
	
	r := rg.Group("/roles", auth.Authenticate())
	{
		r.GET("/", handler.GetRolesTree)
		r.GET("/permissions", handler.GetRolesPermissions)
		r.PUT("/:role_slug/permissions/:permission_name", handler.UpdateRolePermission)
	}
}
