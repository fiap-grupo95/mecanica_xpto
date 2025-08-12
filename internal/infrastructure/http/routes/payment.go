package routes

import (
	"mecanica_xpto/internal/infrastructure/http"

	"github.com/gin-gonic/gin"
)

func addPaymentRoutes(rg *gin.RouterGroup, paymentHandler *http.PaymentHandler) {

	payments := rg.Group(PathPayments)
	{
		payments.GET("/:id", paymentHandler.GetPaymentByID)
		payments.GET("/", paymentHandler.ListPayments)
		payments.POST("/", paymentHandler.CreatePayment)
	}
}
