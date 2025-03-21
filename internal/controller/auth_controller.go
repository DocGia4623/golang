package controller

import (
	"context"
	"log"
	"net/http"
	"testwire/config"
	"testwire/helper"
	"testwire/internal/dto/request"
	"testwire/internal/dto/response"
	"testwire/internal/models"
	"testwire/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthenticationController struct {
	AuthenticationService services.AuthenticationService
	RefreshTokenService   services.RefreshTokenService
}

func NewAuthenticationController(authenticationService services.AuthenticationService, refreshtokenService services.RefreshTokenService) *AuthenticationController {
	return &AuthenticationController{
		AuthenticationService: authenticationService,
		RefreshTokenService:   refreshtokenService,
	}
}

// @Summary Login
// @Description Authenticate user and return access & refresh tokens
// @Tags auth
// @Accept  json
// @Produce  json
// @Param login body request.LoginRequest true "Login details"
// @Success 200 {object} response.LoginResponse
// @Failure 400 {object} response.Response
// @Router /auth/login [post]
func (controller *AuthenticationController) Login(c *gin.Context) {
	LoginRequest := request.LoginRequest{}
	err := c.ShouldBindJSON(&LoginRequest)
	helper.ErrorPanic(err)

	var webResponse response.Response

	refreshToken, accessToken, err := controller.AuthenticationService.Login(LoginRequest.UserName, LoginRequest.Password)
	if err != nil {
		webResponse = response.Response{Code: http.StatusBadRequest, Status: "bad request", Message: "Invalid username or password"}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}
	c.SetCookie("refresh_token", refreshToken, 3600*24*7, "/", "", false, true)
	log.Println("Set refresh token in cookie: " + refreshToken) // ✅ Debug log

	err = controller.RefreshTokenService.SaveToken(models.RefreshToken{
		Token: refreshToken,
	})
	if err != nil {
		webResponse = response.Response{Code: http.StatusBadRequest, Status: "bad request", Message: "Failed saving refresh token"}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}
	webResponse = response.Response{
		Code:    http.StatusOK,
		Status:  "ok",
		Message: "Login success",
		Data:    response.LoginResponse{TokenType: "Bearer Token", RefreshToken: refreshToken, AccessToken: accessToken},
	}
	c.JSON(http.StatusOK, webResponse)
}

// @Summary Register
// @Description Register a new user
// @Tags auth
// @Accept  json
// @Produce  json
// @Param register body request.CreateUserRequest true "User registration details"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /auth/register [post]
func (controller *AuthenticationController) Register(c *gin.Context) {
	UserCreateRequest := request.CreateUserRequest{}
	err := c.ShouldBindJSON(&UserCreateRequest)
	helper.ErrorPanic(err)

	err = controller.AuthenticationService.Register(UserCreateRequest)
	var webResponse response.Response
	if err != nil {
		webResponse = response.Response{Code: http.StatusBadRequest, Status: "bad request", Message: "duplicate username"}
	} else {
		webResponse = response.Response{Code: http.StatusOK, Status: "ok", Message: "Register success"}
	}
	c.JSON(http.StatusOK, webResponse)
}

// @Summary Refresh Token
// @Description Refresh access token using refresh token from cookie
// @Tags auth
// @Accept  json
// @Produce  json
// @Success 200 {object} response.LoginResponse
// @Failure 400 {object} response.Response
// @Router /auth/refresh [post]
func (controller *AuthenticationController) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: http.StatusUnauthorized, Status: "unauthorized", Message: "Refresh token is missing"})
		return
	}
	log.Println("refresh Token: " + refreshToken)
	config, _ := config.LoadConfig()
	newAccessToken, newRefreshToken, err := controller.RefreshTokenService.RefreshToken(refreshToken, config.RefreshTokenSecret)
	c.SetCookie("refresh_token", newRefreshToken, 3600*24*7, "/", "", false, true)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: http.StatusBadRequest, Status: "bad request", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.Response{Code: http.StatusOK, Status: "ok", Message: "Refresh token success", Data: response.LoginResponse{TokenType: "Bearer Token", RefreshToken: newRefreshToken, AccessToken: newAccessToken}})
}

// @Summary Logout
// @Description Logout user by invalidating refresh token
// @Tags auth
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Access token"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /auth/logout [post]
func (controller *AuthenticationController) Logout(c *gin.Context) {
	var webResponse response.Response
	token := c.GetHeader("Authorization")
	if token == "" {
		webResponse = response.Response{
			Code:    http.StatusBadRequest,
			Status:  "bad request",
			Message: "Token is required",
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		webResponse = response.Response{
			Code:    http.StatusBadRequest,
			Status:  "bad request",
			Message: "Refresh token not in cookie",
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}
	err = controller.AuthenticationService.Logout(context.Background(), refreshToken, token)
	if err != nil {
		webResponse = response.Response{
			Code:    http.StatusBadRequest,
			Status:  "bad request",
			Message: "Error during logout: " + err.Error(),
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
	webResponse = response.Response{
		Code:    http.StatusOK,
		Status:  "ok",
		Message: "Logout success ",
	}
	c.JSON(http.StatusOK, webResponse)
}
