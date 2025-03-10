package route

import (
	"golang-template/configs"
	"golang-template/infrastructure/logger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// ! DISABLE GENERATE DOCUMENTATION THROUGH THIS COMMENT IT WORSE
func RegisterSwaggerRoute(router *gin.Engine, cfg *configs.Config, log logger.Logger) {
	if cfg.Environment != "production" {
		log.Info("Enabling Swagger documentation")

		swaggerConfig := ginSwagger.URL("/swagger/doc.json")

		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerConfig))

		router.GET("/docs", func(c *gin.Context) {
			c.Redirect(301, "/swagger/index.html")
		})

		log.Info("Swagger UI available at /swagger/index.html")
	}
}
