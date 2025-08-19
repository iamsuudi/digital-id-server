package city

import (
	"digital-id-server/internal/auth"
	"digital-id-server/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(rg *gin.RouterGroup, db *pgxpool.Pool, q *repository.Queries) {
	service := NewService(db, q)
	handler := NewHandler(service)

	r := rg.Group("/cities", auth.Authenticate())
	{
		r.GET("/", handler.GetCities)
		r.POST("/", handler.CreateCity)
		r.GET("/:id", handler.GetCity)
		r.GET("/:id/audit", handler.GetCityAudit)
		r.PUT("/:id", handler.UpdateCityInfo)
		r.DELETE("/:id", handler.DeleteCity)
		r.GET("/:id/subcities", handler.GetSubCities)
		r.POST("/:id/staff", handler.AddStaff)
		r.DELETE("/:id/staff", handler.RemoveStaff)
	}
}
