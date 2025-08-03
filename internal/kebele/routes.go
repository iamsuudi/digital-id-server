package kebele

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"digital-id-server/internal/auth"
	"digital-id-server/internal/repository"
)

func RegisterRoutes(rg *gin.RouterGroup, db *pgxpool.Pool, q *repository.Queries) {
	service := NewService(db, q)
	handler := NewHandler(service)

	r := rg.Group("/kebeles", auth.Authenticate())
	{
		r.POST("/create", handler.CreateKebele)
		r.GET("/:id", handler.GetKebele)
		r.GET("/:id/encoders", handler.GetKebeleEncoders)
		r.GET("/:id/cashiers", handler.GetKebeleCashiers)
		r.PUT("/:id", handler.UpdateKebele)
		r.GET("/", handler.GetKebeles)
	}
}
