package main

import (
	"github.com/gin-gonic/gin"
	"github.com/iamsuudi/digital-id-server/db"
	"github.com/iamsuudi/digital-id-server/internal/config"
	"github.com/iamsuudi/digital-id-server/internal/routes"
)

func main() {
	cfg := config.Load()
	conn := db.Connect(cfg.DB)
	defer conn.Close()

	r := gin.Default()

	// Register all routes
	routes.Setup(r, conn)

	r.Run(":8080") // Start server
}
