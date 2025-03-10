package dto

import "time"

type HealthResponse struct {
	// Overall health status: healthy, degraded, or unhealthy
	Status string `json:"status" example:"healthy"`
	// Current environment (development, production, etc.)
	Environment string `json:"environment" example:"development"`
	// API version number
	Version string `json:"version" example:"1.0.0"`
	// Current server time in UTC
	Timestamp time.Time `json:"timestamp"`
	// Status of individual services
	Services map[string]Status `json:"services"`
	// Server uptime since start
	Uptime string `json:"uptime" example:"1h23m45s"`
}

// Status represents the status of a service
type Status struct {
	// Service status: ok, degraded, down
	Status string `json:"status" example:"ok"`
	// Optional message with additional information
	Message string `json:"message,omitempty" example:"Running normally"`
}
