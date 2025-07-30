package city

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"digital-id-server/internal/auth"
	"digital-id-server/internal/repository"
)

func RegisterRoutes(rg *gin.RouterGroup, db *pgxpool.Pool, q *repository.Queries) {
	service := NewService(db, q)
	handler := NewHandler(service)

	r := rg.Group("/cities", auth.Authenticate())
	{
		r.POST("/create", handler.CreateCity)
		r.GET("/:id", handler.GetCity)
		r.PUT("/:id", handler.UpdateCity)
		r.GET("/", handler.GetCities)
	}
}
