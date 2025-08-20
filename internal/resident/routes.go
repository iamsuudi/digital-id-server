package resident

import (
	"digital-id-server/internal/auth"
	"digital-id-server/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(rg *gin.RouterGroup, db *pgxpool.Pool, q *repository.Queries) {
	service := NewService(db, q)
	handler := NewHandler(service)

	r := rg.Group("/residents", auth.Authenticate())
	{
		r.GET("/", handler.GetResidents)
		r.GET("/unpaid", handler.GetUnpaidResidents)
		r.GET("/unverified", handler.GetUnverifiedResidents)
		r.POST("/", handler.RegisterResident)
		r.GET("/:id", handler.GetResident)
		r.GET("/:id/audit", handler.GetResidentAudit)
		r.GET("/:id/card", handler.GetIDCard)
		r.GET("/:id/payment", handler.GetResidentPayment)
		r.GET("/:id/documents", handler.GetResidentDocuments)
		r.PUT("/:id/payment", handler.UpdatePaymentInfo)
		r.PUT("/:id/document", handler.UpdateDocumentInfo)
		r.PUT("/:id/personal", handler.UpdatePersonalInfo)
		r.GET("/:id/address", handler.GetAddressInfo)
		r.PUT("/:id/address", handler.UpdateAddressInfo)
		r.PUT("/:id/additional", handler.UpdateAdditionalInfo)
		r.PUT("/:id/emergency", handler.UpdateEmergencyInfo)
		r.PUT("/:id/employment", handler.UpdateEmploymentInfo)
		r.PUT("/:id/biometric", handler.UpdateBiometricInfo)
		r.PUT("/:id/documents", handler.ReplaceDocuments)
	}
}
