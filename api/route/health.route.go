package route

import (
	"golang-template/app/module/health/handler"
	"golang-template/app/module/health/service"
	"golang-template/configs"
	"golang-template/infrastructure/firebase"
	"golang-template/infrastructure/logger"

	"github.com/gin-gonic/gin"
)

func RegisterHealthRoute(router *gin.RouterGroup, cfg *configs.Config, log logger.Logger) {
	fbClient, _ := firebase.Initialize(cfg, log)
	healthService := service.NewHealthService(cfg, fbClient)
	healthHandler := handler.NewHealthHandler(log, healthService)

	router.GET("/health", healthHandler.Check)
}
