package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/iamsuudi/digital-id-server/database/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(rg *gin.RouterGroup, dbConn *pgxpool.Pool, dbQueries *sqlc.Queries) {
	service := NewService(dbConn, dbQueries)
	handler := NewHandler(service)

	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/login", handler.Login)
		authGroup.POST("/logout", handler.Logout)
		authGroup.POST("/register", handler.RegisterUser)
		authGroup.POST("/refresh", handler.RefreshToken)
	}
}
