package routes

import (
	"restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func TableRoutes(router *gin.Engine) {
	router.GET("/tables", controllers.GetTable())
	router.GET("/tables/:table_id", controllers.GetTableByID())
	router.POST("/tables", controllers.CreateTable())
	router.PATCH("/tables/:table_id", controllers.UpdateTable())
}
