package main

import (
	"fmt"
	"log"
	"os"

	"go-web-service/cmd/server/docs"
	productController "go-web-service/internal/products/controller"
	"go-web-service/internal/products/infra"
	productRepository "go-web-service/internal/products/repository/mariadb"
	productService "go-web-service/internal/products/service"

	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(200, gin.H{"message": "oi"})
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		//

		c.Next()
	}
}

// @title MALI Bootcamp API
// @version 1.0
// @description This API Handle MELI Products.
// @termoOfService http://developers.mercadolibre.com.ar/es_ar/terminos-y-condiciones

// @contact.name API Support
// @contact.url http://developers.mercadolibre.com.ar/support

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/License-2.0.html
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	router := gin.Default()

	router.Use(corsMiddleware())

	rep := productRepository.NewMariaDBRepository(infra.Connect())
	service := productService.NewService(rep)

	productController.NewProduto(router, service)

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", os.Getenv("HOST_SERVER"), os.Getenv("PORT_SERVER"))

	docs.SwaggerInfo.BasePath = "/api/v1"

	router.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":" + os.Getenv("PORT_SERVER"))
}
