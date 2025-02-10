package main

import (
	"golang-mongodb/config"
	"golang-mongodb/routes"

	_ "golang-mongodb/docs" // Import dokumentasi Swagger

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Golang JWT Authentication API
// @version 1.0
// @description API untuk autentikasi dengan JWT menggunakan MongoDB
// @host localhost:8080
// @BasePath /
func main() {
	config.ConnectDB()

	r := gin.Default()

	// Tambahkan route Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Routes
	routes.AuthRoutes(r)

	r.Run(":8080")
}
