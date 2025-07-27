package user

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"digital-id-server/internal/auth"
	"digital-id-server/internal/repository"
)

func RegisterRoutes(rg *gin.RouterGroup, dbConn *pgxpool.Pool, dbQueries *repository.Queries) {
	service := NewService(dbConn, dbQueries)
	handler := NewHandler(service)

	usersGroup := rg.Group("/users", auth.Authenticate())
	{
		usersGroup.GET("/", auth.RolePrefixRedirect("superadmin"), handler.GetAll)
		usersGroup.GET("/role", handler.GetByRole)
		usersGroup.GET("/one", handler.GetUser)
	}

	superadmin := rg.Group("/superadmin", auth.Authenticate(), auth.Authorize("superadmin"))
	{
		superadmin.GET("/users", handler.GetAllForSuperadmin)
	}
}
