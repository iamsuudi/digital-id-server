package main

import (
	"github.com/gin-gonic/gin"
	"github.com/iamsuudi/digital-id-server/database"
	"github.com/iamsuudi/digital-id-server/database/sqlc"
	"github.com/iamsuudi/digital-id-server/internal/resident"
	"github.com/iamsuudi/digital-id-server/shared/config"
)

func main() {
	cfg := config.Load()
	dbConn := database.Connect(cfg.DB)
	defer dbConn.Close()

	r := gin.Default()

	// Top-level API group
	api := r.Group("/api")

	// Version 1 group
	v1 := api.Group("/v1")

	dbQueries := sqlc.New(dbConn)

	resident.RegisterRoutes(v1, dbConn, dbQueries)

	r.Run(":8080")
}
