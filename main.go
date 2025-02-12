package main

import (
	"fmt"
	"golang-mongodb/config"
	"golang-mongodb/routes"
	"os"

	_ "golang-mongodb/docs" // Import dokumentasi Swagger

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Golang JWT Authentication API
// @version 1.0
// @description API untuk MongoDB JWT User
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	var jwtKey = []byte(os.Getenv("SECRET_KEY"))

	config.ConnectDB()

	r := gin.Default()

	// Routes
	routes.SetupRoutes(r, jwtKey)

	// Tambahkan route Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}
