package routes

import (
	"testwire/internal/constant"
	"testwire/internal/controller"
	"testwire/internal/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoute(userController controller.UserController, mw *middleware.Middleware, router *gin.Engine) {
	userRoutes := router.Group("/user")
	{
		userRoutes.POST("/addrole", userController.AddRoleForUser)
		userRoutes.GET("", mw.AuthorizeRole(constant.ReadUser), userController.FindAllUser)
	}
}
