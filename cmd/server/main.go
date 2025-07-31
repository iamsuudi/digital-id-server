package main

import (
	"digital-id-server/database"
	"digital-id-server/internal/auth"
	"digital-id-server/internal/cache"
	"digital-id-server/internal/city"
	"digital-id-server/internal/kebele"
	"digital-id-server/internal/permission"
	"digital-id-server/internal/repository"
	"digital-id-server/internal/role"
	"digital-id-server/internal/subcity"

	"github.com/gin-gonic/gin"

	"digital-id-server/internal/user"
	"digital-id-server/shared/config"
)

func main() {
	// Load configurations like environment variables
	config.Load()

	// Connect to database
	db := database.Connect()
	defer db.Close()

	gin.DisableConsoleColor()

	r := gin.Default()

	q := repository.New(db)
	uService := user.NewService(db, q)
	aService := auth.NewService(db, q)
	cache := cache.New(aService, uService)

	auth.RegisterRoutes(r.Group("/api/v1/auth"), db, q)
	city.RegisterRoutes(r.Group("/api/v1/"), db, q)
	user.RegisterRoutes(r.Group("/api/v1"), db, q, cache)
	role.RegisterRoutes(r.Group("/api/v1"), db, q)
	permission.RegisterRoutes(r.Group("/api/v1"), db, q, cache)
	subcity.RegisterRoutes(r.Group("/api/v1"), db, q)
	kebele.RegisterRoutes(r.Group("/api/v1"), db, q)

	r.Run(":8080")
}
