package service

import (
	"context"
	"os"
	"time"

	"golang-template/app/module/health/dto"
	"golang-template/configs"
	"golang-template/infrastructure/firebase"
)

type HealthService interface {
	Check(ctx context.Context) (*dto.HealthResponse, error)
}

type healthServiceImpl struct {
	config     *configs.Config
	firebase   *firebase.Client
	startTime  time.Time
	appVersion string
}

func NewHealthService(config *configs.Config, firebaseClient *firebase.Client) HealthService {
	return &healthServiceImpl{
		config:    config,
		firebase:  firebaseClient,
		startTime: time.Now(),
	}
}

func (h *healthServiceImpl) Check(ctx context.Context) (*dto.HealthResponse, error) {
	firebaseStatus := dto.Status{Status: "ok"}

	if h.firebase != nil && h.firebase.App != nil {
		if h.firebase.Firestore != nil {
			_, err := h.firebase.Firestore.Collections(ctx).GetAll()
			if err != nil {
				firebaseStatus = dto.Status{
					Status:  "degraded",
					Message: "Firestore connection issue",
				}
			}
		} else {
			firebaseStatus = dto.Status{
				Status:  "degraded",
				Message: "Firestore client not initialized",
			}
		}
	} else {
		firebaseStatus = dto.Status{
			Status:  "unknown",
			Message: "Firebase not configured",
		}
	}

	hostname, _ := os.Hostname()

	uptime := time.Since(h.startTime).String()

	services := map[string]dto.Status{
		"firebase":            firebaseStatus,
		"golang-template/api": {Status: "ok"},
		"system": {
			Status:  "ok",
			Message: hostname,
		},
	}

	status := "healthy"
	for _, s := range services {
		if s.Status == "degraded" {
			status = "degraded"
		} else if s.Status == "down" {
			status = "unhealthy"
			break
		}
	}

	return &dto.HealthResponse{
		Status:      status,
		Environment: h.config.Environment,
		Version:     h.appVersion,
		Timestamp:   time.Now().UTC(),
		Services:    services,
		Uptime:      uptime,
	}, nil
}
