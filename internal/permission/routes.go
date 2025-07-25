package permission

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"digital-id-server/internal/auth"
	"digital-id-server/internal/repository"
)

func RegisterRoutes(rg *gin.RouterGroup, dbConn *pgxpool.Pool, dbQueries *repository.Queries) {
	service := NewService(dbConn, dbQueries)
	handler := NewHandler(service)
	
	usersGroup := rg.Group("/permissions", auth.Authenticate())
	{
		usersGroup.GET("/", handler.GetAllPermissions)
		usersGroup.GET("/assignable", handler.GetAssignablePermissions)
	}
}
