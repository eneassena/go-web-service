package main

import (
	"fmt"
	"log"
	"os"

	"web-service-gin/cmd/server/docs"
	productController "web-service-gin/internal/products/controller"
	"web-service-gin/internal/products/infra"
	productRepository "web-service-gin/internal/products/repository/mariadb"
	productService "web-service-gin/internal/products/service"

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

	rep := productRepository.NewMariaDBRepository(infra.Connect())
	service := productService.NewService(rep)

	productController.NewProduto(router, service)

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", os.Getenv("HOST_SERVER"), os.Getenv("PORT_SERVER"))

	docs.SwaggerInfo.BasePath = "/api/v1"

	router.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":" + os.Getenv("PORT_SERVER"))
}
