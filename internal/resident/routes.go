package resident

import (
	"digital-id-server/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(rg *gin.RouterGroup, db *pgxpool.Pool, q *repository.Queries) {
	service := NewService(db, q)
	handler := NewHandler(service)

	r := rg.Group("/residents")
	{
		r.GET("/", handler.GetResidents)
		r.GET("/unpaid", handler.GetUnpaidResidents)
		r.GET("/unverified", handler.GetUnverifiedResidents)
		r.POST("/", handler.RegisterResident)
		r.GET("/:id", handler.GetResident)
		r.GET("/:id/payment", handler.GetResidentPayment)
		r.GET("/:id/documents", handler.GetResidentDocuments)
		r.PUT("/:id/payment", handler.UpdatePaymentInfo)
		r.PUT("/:id/document", handler.UpdateDocumentInfo)
	}
}
