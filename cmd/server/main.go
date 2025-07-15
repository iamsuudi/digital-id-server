package main

import (
	"github.com/gin-gonic/gin"
	"github.com/iamsuudi/digital-id-server/database"
	"github.com/iamsuudi/digital-id-server/internal/auth"
	"github.com/iamsuudi/digital-id-server/internal/city"
	"github.com/iamsuudi/digital-id-server/internal/repository"
	"github.com/iamsuudi/digital-id-server/internal/resident"
	"github.com/iamsuudi/digital-id-server/shared/config"
)

func main() {
	config.Load()
	dbConn := database.Connect()
	defer dbConn.Close()
	gin.DisableConsoleColor()
	r := gin.Default()

	// Top-level API group
	api := r.Group("/api")

	// Version 1 group
	v1 := api.Group("/v1")

	dbQueries := repository.New(dbConn)

	auth.RegisterRoutes(v1, dbConn, dbQueries)
	resident.RegisterRoutes(v1, dbConn, dbQueries)
	city.RegisterRoutes(v1, dbConn, dbQueries)

	r.Run(":8080")
}
