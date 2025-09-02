package main

import (
	"github.com/gin-gonic/gin"

	"digital-id-server/database"
	"digital-id-server/internal/analytics"
	"digital-id-server/internal/audit"
	"digital-id-server/internal/auth"
	"digital-id-server/internal/cache"
	"digital-id-server/internal/city"
	"digital-id-server/internal/kebele"
	"digital-id-server/internal/permission"
	"digital-id-server/internal/repository"
	"digital-id-server/internal/resident"
	"digital-id-server/internal/role"
	"digital-id-server/internal/subcity"
	"digital-id-server/internal/user"
	"digital-id-server/shared/config"
	"digital-id-server/shared/email"
)

func main() {
	config.Load()

	db := database.Connect()
	defer db.Close()

	gin.DisableConsoleColor()

	r := gin.Default()

	r.Static("/assets", "./uploads")

	q := repository.New(db)
	uService := user.NewService(db, q)
	aService := auth.NewService(db, q)
	cache := cache.New(aService, uService)

	e, err := email.DefaultService()
	if err != nil {
		panic(err)
	}

	// eService.Send(&email.EmailParams{
	// 	To:      []string{"suudiabdulfetah@gmail.com"},
	// 	Subject: "Welcome",
	// 	Text:    "Welcome to Oict Digital ID",
	// })

	auth.RegisterRoutes(r.Group("/api/v1/auth"), db, q, e)
	city.RegisterRoutes(r.Group("/api/v1/"), db, q)
	user.RegisterRoutes(r.Group("/api/v1"), db, q, cache)
	role.RegisterRoutes(r.Group("/api/v1"), db, q)
	permission.RegisterRoutes(r.Group("/api/v1"), db, q, cache)
	subcity.RegisterRoutes(r.Group("/api/v1"), db, q)
	kebele.RegisterRoutes(r.Group("/api/v1"), db, q)
	audit.RegisterRoutes(r.Group("/api/v1"), db, q, cache)
	resident.RegisterRoutes(r.Group("/api/v1"), db, q)
	analytics.RegisterRoutes(r.Group("/api/v1"), db, q, cache)

	r.Run(":8080")
}
