package main

import (
	"github.com/gin-gonic/gin"
	"digital-id-server/database"
	"digital-id-server/internal/auth"
	"digital-id-server/internal/city"
	"digital-id-server/internal/permission"
	"digital-id-server/internal/repository"
	"digital-id-server/internal/role"

	// "digital-id-server/internal/resident"
	"digital-id-server/internal/user"
	"digital-id-server/shared/config"
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
	// resident.RegisterRoutes(v1, dbConn, dbQueries)
	city.RegisterRoutes(v1, dbConn, dbQueries)
	user.RegisterRoutes(v1, dbConn, dbQueries)
	role.RegisterRoutes(v1, dbConn, dbQueries)
	permission.RegisterRoutes(v1, dbConn, dbQueries)

	r.Run(":8080")
}
