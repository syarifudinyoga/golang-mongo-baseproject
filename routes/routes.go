package routes

import (
	"golang-mongodb/handlers"
	"golang-mongodb/utils"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, jwtSecret []byte) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	users := router.Group("/users", utils.AuthMiddleware(jwtSecret))
	{
		users.GET("/", handlers.GetAllUsers)
		users.PUT("/:id", handlers.UpdateUser)
		users.DELETE("/:id", handlers.DeleteUser)
	}
}
