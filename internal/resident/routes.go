package resident

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"digital-id-server/internal/repository"
)

func RegisterRoutes(rg *gin.RouterGroup, db *pgxpool.Pool, q *repository.Queries) {
	service := NewService(db, q)
	handler := NewHandler(service)

	r := rg.Group("/residents")
	{
		r.GET("/", handler.GetResidents)
		r.POST("/", handler.RegisterResident)
		r.GET("/:id", handler.GetResident)
	}
}
