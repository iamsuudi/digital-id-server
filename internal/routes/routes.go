package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/iamsuudi/digital-id-server/internal/user"
	generated "github.com/iamsuudi/digital-id-server/db/generated"
)

func Setup(r *gin.Engine, db *pgxpool.Pool) {
	api := r.Group("/api")

	queries := generated.New(db)

	userHandler := user.NewHandler(queries)
	api.GET("/users/:id", userHandler.GetUserByID)
	api.POST("/users", userHandler.CreateUser)

	// Add more module routes here
}
