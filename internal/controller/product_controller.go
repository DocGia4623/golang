package controller

import (
	"net/http"
	"testwire/internal/dto/response"
	"testwire/internal/models"
	"testwire/internal/services"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	ProductService services.ProductService
}

func NewProductController(productService services.ProductService) *ProductController {
	return &ProductController{ProductService: productService}
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
	c.JSON(http.StatusInternalServerError, webResponse)
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
