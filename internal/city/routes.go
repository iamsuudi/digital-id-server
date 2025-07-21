package city

import (
	"github.com/gin-gonic/gin"
	"github.com/iamsuudi/digital-id-server/internal/auth"
	"github.com/iamsuudi/digital-id-server/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(rg *gin.RouterGroup, dbConn *pgxpool.Pool, dbQueries *repository.Queries) {
	service := NewService(dbConn, dbQueries)
	handler := NewHandler(service)

	residentGroup := rg.Group("/cities", auth.Authenticate())
	{
		residentGroup.POST("/create", handler.CreateCity)
		residentGroup.GET("/:id", handler.GetCity)
		residentGroup.GET("/", handler.GetAll)
	}
}
