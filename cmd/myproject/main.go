package main

import (
	"backend/config"
	"backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Koneksi ke database
	config.ConnectDB()
	defer config.DB.Close()

	// Setup router Gin
	router := gin.Default()

	// Setup Routes
	routes.SetupRoutes(router, config.DB)

	// Jalankan server
	router.Run(":8080")
}
