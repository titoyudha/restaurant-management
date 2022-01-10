package routes

import (
	"restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func MenuRoutes(router *gin.Engine) {
	router.GET("/menus", controllers.GetMenu())
	router.GET("/menus/:menu_id", controllers.GetMenuByID())
	router.POST("/menus", controllers.CreateMenu())
	router.PATCH("/menus/:menu_id", controllers.UpdateMenu())
}
