package permission

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"digital-id-server/internal/auth"
	"digital-id-server/internal/repository"
)

func RegisterRoutes(rg *gin.RouterGroup, db *pgxpool.Pool, q *repository.Queries) {
	service := NewService(db, q)
	handler := NewHandler(service)
	
	r := rg.Group("/permissions", auth.Authenticate())
	{
		r.GET("/", handler.GetAllPermissions)
		r.GET("/assignable", handler.GetAssignablePermissions)
	}
}
