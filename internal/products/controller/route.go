package controller

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(route *gin.Engine, productController *ProdutoController) {
	group := route.Group("/api/v1")
	group.GET("/products", productController.GetAll())
	group.GET("/products/:id", productController.GetOne())
	group.POST("/products", productController.Store())
	group.PUT("/products/:id", productController.Update())
	group.PATCH("/products/:id", productController.UpdateName())
	group.DELETE("/products/:id", productController.Delete())
}
