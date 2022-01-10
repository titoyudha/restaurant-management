package routes

import (
	"restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	router.GET("/users", controllers.GetUser())
	router.GET("/user/:user_id", controllers.GetUserById())
	router.POST("/users/signup", controllers.SignUp())
	router.POST("/users/login", controllers.LogIn())
}
