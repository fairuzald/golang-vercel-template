package configs

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Application
	AppName     string
	Environment string
	Port        int
	Debug       bool

	// Firebase
	FirebaseProjectID   string
	FirebaseCredentials string

	// API Rate Limiting
	RateLimitRequests int
	RateLimitDuration time.Duration

	// CORS
	CORSAllowedOrigins []string
	CORSAllowedMethods []string
	CORSAllowedHeaders []string
	CORSExposedHeaders []string
	CORSMaxAge         time.Duration

	// Security
	AuthTokenExpiry time.Duration
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		// Application
		AppName:     getEnv("APP_NAME", "golang-template"),
		Environment: getEnv("APP_ENV", "development"),
		Port:        getEnvAsInt("APP_PORT", 8080),
		Debug:       getEnvAsBool("APP_DEBUG", true),

		// Firebase
		FirebaseProjectID:   getEnv("FIREBASE_PROJECT_ID", ""),
		FirebaseCredentials: getEnv("FIREBASE_SERVICE_ACCOUNT", "./credentials/firebase-service-account.json"),

		// API Rate Limiting
		RateLimitRequests: getEnvAsInt("RATE_LIMIT_REQUESTS", 100),
		RateLimitDuration: getEnvAsDuration("RATE_LIMIT_DURATION", 1*time.Minute),

		// CORS
		CORSAllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", "*"),
		CORSAllowedMethods: getEnvAsSlice("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS"),
		CORSAllowedHeaders: getEnvAsSlice("CORS_ALLOWED_HEADERS", "Authorization,Content-Type,X-Requested-With"),
		CORSExposedHeaders: getEnvAsSlice("CORS_EXPOSED_HEADERS", "Content-Length"),
		CORSMaxAge:         getEnvAsDuration("CORS_MAX_AGE", 12*time.Hour),

		// Security
		AuthTokenExpiry: getEnvAsDuration("AUTH_TOKEN_EXPIRY", 24*time.Hour),
	}
}

func (c *Config) GetFirebaseCredentialsPath() string {
	if filepath.IsAbs(c.FirebaseCredentials) {
		return c.FirebaseCredentials
	}

	dir, err := os.Getwd()
	if err != nil {
		return c.FirebaseCredentials
	}

	// Handle vercel deployment issues path
	if strings.Contains(dir, ".vercel/cache/go/") {
		projectRoot := filepath.Join(dir, "../../../..")
		absPath := filepath.Join(projectRoot, c.FirebaseCredentials)
		if _, err := os.Stat(absPath); err == nil {
			return absPath
		}
	}

	return filepath.Join(dir, c.FirebaseCredentials)
}

// Helper functions to get environment and parse values
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue string) []string {
	valueStr := getEnv(key, defaultValue)
	if valueStr == "" {
		return []string{}
	}
	return strings.Split(valueStr, ",")
}
