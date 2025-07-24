package permission

import (
	"github.com/gin-gonic/gin"
	"github.com/iamsuudi/digital-id-server/internal/auth"
	"github.com/iamsuudi/digital-id-server/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
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
