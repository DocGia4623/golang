package controller

import (
	"net/http"
	"testwire/internal/dto/response"
	"testwire/internal/models"
	"testwire/internal/services"
	"testwire/logs"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ProductController struct {
	ProductService services.ProductService
	UserService    services.UserSerive
}

func NewProductController(productService services.ProductService, userService services.UserSerive) *ProductController {
	return &ProductController{ProductService: productService, UserService: userService}
}

func (controller *ProductController) GetSearchHistory(c *gin.Context) {
	// Query lấy lịch sử tìm kiếm từ Elasticsearch
	// query := `{
	// 	"query": {
	// 		"match": { "message": "User search successful" }
	// 	},
	// 	"sort": [ { "timestamp": { "order": "desc" } } ]
	// }`
}
func (controller *ProductController) FindByName(c *gin.Context) {
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
	name := c.PostForm("name")
	if name == "" {
		webResponse = response.Response{
			Code:    http.StatusBadRequest,
			Status:  "fail",
			Message: "name is required",
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}
	products, err := controller.ProductService.FindByName(name)
	if err != nil {
		webResponse = response.Response{
			Code:    http.StatusInternalServerError,
			Status:  "fail",
			Message: "Server error",
		}
		c.JSON(http.StatusInternalServerError, webResponse)
		return
	}
	userId, err := controller.UserService.GetUserID(token)
	if err != nil {
		webResponse = response.Response{
			Code:    http.StatusInternalServerError,
			Status:  "fail",
			Message: "Invalid token: " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, webResponse)
		return
	}
	logs.Init()
	logs.Logger.WithFields(logrus.Fields{
		"userId":    userId,
		"Data":      products,
		"message":   "User search successful for name: " + name,
		"level":     "info",
		"timestamp": time.Now().Format(time.RFC3339),
	}).Info("User search event")
	logs.CloseLogFile()

	// ✅ Ghi lịch sử tìm kiếm vào file log (Logstash sẽ tự động thu thập)

	webResponse = response.Response{
		Code:    http.StatusOK,
		Status:  "ok",
		Message: "Showing product:",
		Data:    products,
	}
	c.JSON(http.StatusOK, webResponse)
}

func (controller *ProductController) CreateProduct(c *gin.Context) {
	var webResponse response.Response
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		webResponse = response.Response{
			Code:    http.StatusBadRequest,
			Status:  "fail",
			Message: "invalid request format",
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}
	err := controller.ProductService.SaveProduct(product)
	if err != nil {
		webResponse = response.Response{
			Code:    http.StatusInternalServerError,
			Status:  "fail",
			Message: "Server error",
		}
		c.JSON(http.StatusInternalServerError, webResponse)
		return
	}
	webResponse = response.Response{
		Code:    http.StatusOK,
		Status:  "ok",
		Message: "Created product",
		Data:    product,
	}
	c.JSON(http.StatusOK, webResponse)
}

func (controller *ProductController) DeleteProduct(c *gin.Context) {
	var webResponse response.Response
	var name string
	err := c.ShouldBindJSON(&name)
	if err != nil {
		webResponse = response.Response{
			Code:    http.StatusInternalServerError,
			Status:  "fail",
			Message: "invalid request format",
		}
		c.JSON(http.StatusInternalServerError, webResponse)
		return
	}
	err = controller.ProductService.DeleteProduct(name)
	if err != nil {
		webResponse = response.Response{
			Code:    http.StatusInternalServerError,
			Status:  "fail",
			Message: "Server error",
		}
		c.JSON(http.StatusInternalServerError, webResponse)
		return
	}
	webResponse = response.Response{
		Code:    http.StatusOK,
		Status:  "ok",
		Message: "Deleted product",
	}
	c.JSON(http.StatusInternalServerError, webResponse)
}

func (controller *ProductController) GetAll(c *gin.Context) {
	var webResponse response.Response
	products, err := controller.ProductService.GetAll()
	if err != nil {
		webResponse = response.Response{
			Code:    http.StatusInternalServerError,
			Status:  "fail",
			Message: "Server error: %w" + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, webResponse)
		return
	}
	webResponse = response.Response{
		Code:    http.StatusOK,
		Status:  "ok",
		Message: "Showing product:",
		Data:    products,
	}
	c.JSON(http.StatusInternalServerError, webResponse)
}
