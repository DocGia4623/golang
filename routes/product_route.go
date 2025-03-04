package routes

import (
	"testwire/internal/controller"
	"testwire/internal/middleware"

	"github.com/gin-gonic/gin"
)

func ProductRoute(productController controller.ProductController, mw *middleware.Middleware, router *gin.Engine) {
	productRoutes := router.Group("/product")
	{
		productRoutes.GET("", productController.GetAll)
		productRoutes.POST("/create", productController.CreateProduct)
		productRoutes.POST("/delete", productController.DeleteProduct)
	}
}
