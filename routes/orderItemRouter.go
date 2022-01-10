package routes

import (
	"restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func OrderItemsRoutes(router *gin.Engine) {
	router.GET("/orderItems", controllers.GetOrderItems())
	router.GET("/orderItems/:orderItem_id", controllers.GetOrderItemByID())
	router.GET("/orderItems-order/:order_id", controllers.GetOrderItemsbyOrder())
	router.POST("/orderItems", controllers.CreateOrderItem())
	router.PATCH("/orderItems/:orderItem_id", controllers.UpdateOrderItem())
}
