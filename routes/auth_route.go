package routes

import (
	"testwire/internal/controller"

	"github.com/gin-gonic/gin"
)

// AuthRoute sets up authentication routes
func AuthRoute(authController controller.AuthenticationController, router *gin.Engine) {
	authRoutes := router.Group("/auth")
	{
		// @Summary Login
		// @Description Login
		// @Tags auth
		// @Accept  json
		// @Produce  json
		// @Param login body request.LoginRequest true "Login"
		// @Success 200 {object} response.LoginResponse
		// @Router /auth/login [post]
		authRoutes.POST("/login", authController.Login)

		// @Summary Register
		// @Description Register
		// @Tags auth
		// @Accept  json
		// @Produce  json
		// @Param register body request.CreateUserRequest true "Register"
		// @Success 200 {object} response.Response
		// @Router /auth/register [post]
		authRoutes.POST("/register", authController.Register)

		// @Summary Refresh
		// @Description Refresh Access Token using refresh token
		// @Tags auth
		// @Accept  json
		// @Produce  json
		// @Success 200 {object} response.Response
		// @Router /auth/refresh [post]
		authRoutes.POST("/refresh", authController.RefreshToken)

		// @Summary Logout
		// @Description Logout
		// @Tags auth
		// @Accept  json
		// @Produce  json
		// @Success 200 {object} response.Response
		// @Router /auth/logout [post]
		authRoutes.POST("/logout", authController.Logout)
	}
}
