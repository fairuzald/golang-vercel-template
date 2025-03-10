package handler

import (
	"golang-template/app/module/health/service"
	"golang-template/infrastructure/logger"
	"golang-template/pkg/common/response"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	logger  logger.Logger
	service service.HealthService
}

func NewHealthHandler(logger logger.Logger, service service.HealthService) *HealthHandler {
	return &HealthHandler{
		logger:  logger,
		service: service,
	}
}

func (h *HealthHandler) Check(c *gin.Context) {
	result, err := h.service.Check(c.Request.Context())
	if err != nil {
		h.logger.Error("Health check failed", "error", err)
		response.Error(c, err)
		return
	}

	response.OK(c, result)
}
