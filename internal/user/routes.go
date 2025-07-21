package user

import (
	"github.com/gin-gonic/gin"
	"github.com/iamsuudi/digital-id-server/internal/auth"
	"github.com/iamsuudi/digital-id-server/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(rg *gin.RouterGroup, dbConn *pgxpool.Pool, dbQueries *repository.Queries) {
	service := NewService(dbConn, dbQueries)
	handler := NewHandler(service)

	usersGroup := rg.Group("/users", auth.Authenticate(), auth.RolePrefixRedirect("superadmin"))
	{
		usersGroup.GET("/", handler.GetAll)
		usersGroup.GET("/one", handler.GetUser)
	}

	superadmin := rg.Group("/superadmin", auth.Authenticate(), auth.Authorize("superadmin"))
	{
		superadmin.GET("/users", handler.GetAllForSuperadmin)
	}
}
