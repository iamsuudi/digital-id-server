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
	conn := database.Connect(cfg.DB)
	defer conn.Close()

	r := gin.Default()
	api := r.Group("/api")

	dbQueries := sqlc.New(conn)

	// Setup resident routes, injecting residentService
	residentService := resident.NewService(dbQueries, conn)
	residentHandler := resident.NewHandler(residentService)
	resident.RegisterRoutes(api, residentHandler)

	r.Run(":8080")
}
