package controller

import (
	"net/http"
	"testwire/internal/dto/request"
	"testwire/internal/dto/response"
	"testwire/internal/services"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService services.UserSerive
}

func NewUserController(userService services.UserSerive) *UserController {
	return &UserController{UserService: userService}
}

func (controller *UserController) GetUsesrId(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")
	if header == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Token is required"})
		c.Abort()
		return
	}
	userID, err := controller.UserService.GetUserID(header)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"user_id": userID})
}

func (controller *UserController) AddRoleForUser(c *gin.Context) {
	var webResponse response.Response

	var req request.AddRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		webResponse = response.Response{
			Code:    http.StatusBadRequest,
			Status:  "fail",
			Message: "invalid request format",
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}

	if req.UserID == 0 {
		webResponse = response.Response{
			Code:    http.StatusBadRequest,
			Status:  "fail",
			Message: "user ID is required",
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}

	roles, err := controller.UserService.FindRole(req.Roles)
	if err != nil {
		webResponse = response.Response{
			Code:    http.StatusBadRequest,
			Status:  "fail",
			Message: "error finding roles",
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}
	err = controller.UserService.AddRole(req.UserID, roles)
	if err != nil {
		webResponse = response.Response{
			Code:    http.StatusBadRequest,
			Status:  "fail",
			Message: "Add role err :" + err.Error(),
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}
	webResponse = response.Response{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Add role success",
	}
	c.JSON(http.StatusOK, webResponse)
}

func (controller *UserController) FindAllUser(c *gin.Context) {
	var webResponse response.Response

	users, err := controller.UserService.FindAll()
	if err != nil {
		webResponse = response.Response{
			Code:    http.StatusBadRequest,
			Status:  "fail",
			Message: "error finding user :" + err.Error(),
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}
	webResponse = response.Response{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Listing users",
		Data:    users,
	}
	c.JSON(http.StatusOK, webResponse)
}
