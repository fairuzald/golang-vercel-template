package route

import (
	"net/http"

	"golang-template/configs"
	"golang-template/infrastructure/logger"
	"golang-template/pkg/common/response"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, cfg *configs.Config, log logger.Logger) {

	// API routes group
	apiGroup := router.Group("/api")
	{
		RegisterHealthRoute(apiGroup, cfg, log)
	}

	RegisterSwaggerRoute(router, cfg, log)

	// Home route
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to GolangTemplate API")
	})

	// 404 handler
	router.NoRoute(func(c *gin.Context) {
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
			response.NotFound(c, "API route not found")
			return
		}
		c.String(http.StatusNotFound, "Page not found")
	})

	// 403 method not allowed
	router.NoMethod(func(c *gin.Context) {
		response.ErrorWithCode(c, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
	})
}
