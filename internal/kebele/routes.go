package kebele

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"digital-id-server/internal/auth"
	"digital-id-server/internal/repository"
)

func RegisterRoutes(rg *gin.RouterGroup, dbConn *pgxpool.Pool, dbQueries *repository.Queries) {
	service := NewService(dbConn, dbQueries)
	handler := NewHandler(service)

	group := rg.Group("/kebeles", auth.Authenticate())
	{
		group.POST("/create", handler.CreateKebele)
		group.GET("/:id", handler.GetKebele)
		group.PUT("/:id", handler.UpdateKebele)
		group.GET("/", handler.GetAll)
	}
}
