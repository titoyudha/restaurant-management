package routes

import (
	"restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func OrderRoutes(router *gin.Engine) {
	router.GET("/orders", controllers.GetOrder())
	router.GET("/orders/:order_id", controllers.GetOrderByID())
	router.POST("/orders", controllers.CreateOrder())
	router.PATCH("/orders/:order_id", controllers.UpdateOrder())
}
