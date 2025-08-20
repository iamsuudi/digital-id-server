package subcity

import (
	"digital-id-server/internal/auth"
	"digital-id-server/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(rg *gin.RouterGroup, db *pgxpool.Pool, q *repository.Queries) {
	service := NewService(db, q)
	handler := NewHandler(service)

	r := rg.Group("/subcities", auth.Authenticate())
	{
		r.GET("/", handler.GetSubCities)
		r.POST("/", handler.CreateSubCity)
		r.GET("/:id", handler.GetSubCity)
		r.GET("/:id/audit", handler.GetSubCityAudit)
		r.DELETE("/:id", handler.DeleteSubCity)
		r.PUT("/:id", handler.UpdateSubCityInfo)
		r.GET("/:id/kebeles", handler.GetKebeles)
		r.POST("/:id/staff", handler.AddStaff)
		r.DELETE("/:id/staff", handler.RemoveStaff)
	}
}
