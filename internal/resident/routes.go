package resident

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"digital-id-server/internal/repository"
)

func RegisterRoutes(rg *gin.RouterGroup, dbConn *pgxpool.Pool, dbQueries *repository.Queries) {
	service := NewService(dbConn, dbQueries)
	handler := NewHandler(service)

	residentGroup := rg.Group("/residents")
	{
		residentGroup.POST("/register", handler.RegisterResident)
		residentGroup.GET("/:id", handler.GetResident)
		residentGroup.GET("/", handler.GetAll)
	}
}
