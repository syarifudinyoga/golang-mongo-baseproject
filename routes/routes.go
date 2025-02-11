package routes

import (
	"golang-mongodb/handlers"

	"github.com/gin-gonic/gin"
)

func RoutesGo(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	users := router.Group("/users")
	{
		users.GET("/", handlers.GetAllUsers)
		users.PUT("/:id", handlers.UpdateUser)
		users.DELETE("/:id", handlers.DeleteUser)
	}
}
