package role

import (
	"github.com/gin-gonic/gin"
	"github.com/iamsuudi/digital-id-server/internal/auth"
	"github.com/iamsuudi/digital-id-server/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(rg *gin.RouterGroup, dbConn *pgxpool.Pool, dbQueries *repository.Queries) {
	service := NewService(dbConn, dbQueries)
	handler := NewHandler(service)
	
	usersGroup := rg.Group("/roles", auth.Authenticate())
	{
		usersGroup.GET("/", handler.GetRolesTree)
		usersGroup.GET("/permissions", handler.GetRolesPermissions)
		usersGroup.GET("/rank", handler.GetMyRoleLevelRank)
		usersGroup.PUT("/:role_slug/permissions/:permission_name", handler.UpdateRolePermission)
	}
}
