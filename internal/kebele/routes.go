package kebele

import (
	"digital-id-server/internal/auth"
	"digital-id-server/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(rg *gin.RouterGroup, db *pgxpool.Pool, q *repository.Queries) {
	service := NewService(db, q)
	handler := NewHandler(service)

	r := rg.Group("/kebeles", auth.Authenticate())
	{
		r.GET("/", handler.GetKebeles)
		r.POST("/create", handler.CreateKebele)
		r.GET("/:id", handler.GetKebele)
		r.GET("/:id/encoders", handler.GetKebeleEncoders)
		r.GET("/:id/cashiers", handler.GetKebeleCashiers)
		r.PUT("/:id", handler.UpdateKebeleInfo)
		r.POST("/:id/staff", handler.AddStaff)
		r.DELETE("/:id/staff", handler.RemoveStaff)
	}
}
