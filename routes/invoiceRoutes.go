package routes

import (
	"restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func InvoiceRoutes(router *gin.Engine) {
	router.GET("/invoices", controllers.GetInvoice())
	router.GET("/invoices/:invoice_id", controllers.GetInvoiceByID())
	router.POST("/invoice", controllers.CreateInvoice())
	router.PATCH("/invoice/:invoice_id", controllers.UpdateInvoice())
}
