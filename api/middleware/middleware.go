package middleware

import (
	"os"
	"time"

	"golang-template/configs"
	"golang-template/infrastructure/logger"
	"golang-template/pkg/common/response"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func Setup(router *gin.Engine, cfg *configs.Config, log logger.Logger) {
	// logger middleware
	router.Use(ginzap.Ginzap(log.ZapLogger(), time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(log.ZapLogger(), true))

	// CORS middleware
	router.Use(corsMiddleware(cfg))

	// rate limiting middleware
	rateLimiterMiddleware(router, cfg)

	// security headers
	router.Use(securityHeadersMiddleware())

	// request ID middleware
	router.Use(requestIDMiddleware())
}

// corsMiddleware configures CORS
func corsMiddleware(cfg *configs.Config) gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowOrigins = cfg.CORSAllowedOrigins
	config.AllowMethods = cfg.CORSAllowedMethods
	config.AllowHeaders = cfg.CORSAllowedHeaders
	config.ExposeHeaders = cfg.CORSExposedHeaders
	config.MaxAge = cfg.CORSMaxAge
	config.AllowCredentials = true

	return cors.New(config)
}

// Rate limiter middleware
func rateLimiterMiddleware(router *gin.Engine, cfg *configs.Config) {
	rate := limiter.Rate{
		Period: cfg.RateLimitDuration,
		Limit:  int64(cfg.RateLimitRequests),
	}

	store := memory.NewStore()

	instance := limiter.New(store, rate)

	rateLimiter := mgin.NewMiddleware(instance)

	router.Use(func(c *gin.Context) {

		rateLimiter(c)

		//  context was aborted by the rate limiter
		if c.Writer.Status() == 429 {
			c.Abort()
			response.RateLimitExceeded(c)
			return
		}
	})
}

func securityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-XSS-Protection", "1; mode=block")

		if os.Getenv("APP_ENV") == "production" {
			c.Header("X-Frame-Options", "DENY")
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

			c.Header("Content-Security-Policy",
				"default-src 'self'; "+
					"script-src 'self'; "+
					"style-src 'self' 'unsafe-inline' cdnjs.cloudflare.com fonts.googleapis.com; "+
					"font-src 'self' fonts.gstatic.com cdnjs.cloudflare.com; "+
					"img-src 'self' data: images.unsplash.com; "+
					"connect-src 'self'")
		} else {
			c.Header("Content-Security-Policy",
				"default-src * 'unsafe-inline' 'unsafe-eval'; "+
					"img-src * data:; "+
					"font-src * data:; "+
					"connect-src *")
		}

		c.Next()
	}
}

// requestIDMiddleware  request ID to each request
func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("RequestID", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}
