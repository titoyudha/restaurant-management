package main

import (
	"restaurant-management/databases"
	"restaurant-management/middlewares"
	"restaurant-management/routes"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	foodCollection *mongo.Collection = databases.OpenCollection(databases.Client, "food")
)

func main() {
	router := gin.Default()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(middlewares.Authentication())

	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.TableRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemsRoutes(router)
	routes.InvoiceRoutes(router)

	router.Run(":8080")
}
