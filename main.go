package main

import (
	"github.com/gin-gonic/gin"
	"github.com/iamsuudi/digital-id-server/config"
	"github.com/iamsuudi/digital-id-server/models"
	"github.com/iamsuudi/digital-id-server/routes"
)

func main() {
	config.InitDB()
    config.DB.AutoMigrate(&models.User{})
	
	r := gin.Default()

	routes.RegisterRoutes(r)
	r.Run()
}
