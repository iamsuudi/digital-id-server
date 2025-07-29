package main

import (
	"digital-id-server/database"
	"digital-id-server/internal/auth"
	"digital-id-server/internal/city"
	"digital-id-server/internal/kebele"
	"digital-id-server/internal/cache"
	"digital-id-server/internal/permission"
	"digital-id-server/internal/repository"
	"digital-id-server/internal/role"
	"digital-id-server/internal/subcity"

	"github.com/gin-gonic/gin"

	"digital-id-server/internal/user"
	"digital-id-server/shared/config"
)

func main() {
	config.Load()
	db := database.Connect()
	defer db.Close()
	gin.DisableConsoleColor()
	r := gin.Default()

	api := r.Group("/api")
	v1 := api.Group("/v1")

	q := repository.New(db)
	uService := user.NewService(db, q)
	aService := auth.NewService(db, q)
	cache := cache.New(aService, uService)

	auth.RegisterRoutes(v1, db, q)
	city.RegisterRoutes(v1, db, q)
	user.RegisterRoutes(v1, db, q, cache)
	role.RegisterRoutes(v1, db, q)
	permission.RegisterRoutes(v1, db, q)
	subcity.RegisterRoutes(v1, db, q)
	kebele.RegisterRoutes(v1, db, q)

	r.Run(":8080")
}
