package controller

import (
	"net/http"
	"testwire/config"
	"testwire/internal/dto/response"
	"testwire/internal/services"
	"testwire/utils"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	OrderService services.OrderSerivce
}

func NewOrderController(orderService services.OrderSerivce) *OrderController {
	return &OrderController{OrderService: orderService}
}

func (o *OrderController) Order(c *gin.Context) {
	var webResponse response.Response
	authHeader := c.GetHeader("authorization")
	sub, err := utils.ValidateAccessToken(authHeader, config.AppConfig.AccessTokenSecret)
	if err != nil {
		webResponse = response.Response{
			Code:    http.StatusInternalServerError,
			Status:  "fail",
			Message: "err validate token : " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, webResponse)
		return
	}
	userId, ok := sub.(uint)
	if !ok {
		webResponse = response.Response{
			Code:    http.StatusInternalServerError,
			Status:  "fail",
			Message: "error parsing id",
		}
		c.JSON(http.StatusInternalServerError, webResponse)
		return
	}
}
