package subcity

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"digital-id-server/internal/auth"
	"digital-id-server/internal/repository"
)

func RegisterRoutes(rg *gin.RouterGroup, db *pgxpool.Pool, q *repository.Queries) {
	service := NewService(db, q)
	handler := NewHandler(service)

	r := rg.Group("/subcities", auth.Authenticate())
	{
		r.POST("/create", handler.CreateSubCity)
		r.GET("/:id", handler.GetSubCity)
		r.PUT("/:id", handler.UpdateSubCityInfo)
		r.GET("/:id/kebeles", handler.GetKebeles)
		r.POST("/:id/staff", handler.AddStaff)
		r.DELETE("/:id/staff", handler.RemoveStaff)
		r.GET("/", handler.GetSubCities)
	}
}
