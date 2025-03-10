package api

import (
	"net/http"

	"golang-template/api/middleware"
	"golang-template/api/route"
	"golang-template/configs"
	"golang-template/infrastructure/logger"

	"github.com/gin-gonic/gin"
)

// Handler is the entry point for DEPLOYMENT
func Handler(w http.ResponseWriter, r *http.Request) {
	cfg := configs.LoadConfig()
	log := logger.NewLogger(cfg.Debug)
	router := SetupRouter(cfg, log)
	router.ServeHTTP(w, r)
}

func SetupRouter(cfg *configs.Config, log logger.Logger) *gin.Engine {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	middleware.Setup(router, cfg, log)

	route.RegisterRoutes(router, cfg, log)

	return router
}
