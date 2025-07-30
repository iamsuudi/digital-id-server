package auth

import (
	"digital-id-server/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(rg *gin.RouterGroup, db *pgxpool.Pool, q *repository.Queries) {
	service := NewService(db, q)
	handler := NewHandler(service)

	rg.POST("/login", handler.Login)
	rg.POST("/logout", handler.Logout)
	rg.POST("/register", handler.RegisterUser)
	rg.POST("/refresh", handler.RefreshToken)
	rg.GET("/me", Authenticate(), handler.Me)
}
