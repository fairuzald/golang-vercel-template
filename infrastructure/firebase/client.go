package firebase

import (
	"context"
	"os"
	"strings"
	"sync"

	"golang-template/configs"
	"golang-template/infrastructure/logger"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/storage"
	"google.golang.org/api/option"
)

type Client struct {
	App       *firebase.App
	Auth      *auth.Client
	Firestore *firestore.Client
	Storage   *storage.Client
	Config    *configs.Config
	Logger    logger.Logger
}

var (
	instance *Client
	once     sync.Once
)

// Initialize sets up the Firebase client
func Initialize(cfg *configs.Config, log logger.Logger) (*Client, error) {
	var err error

	once.Do(func() {
		instance = &Client{
			Config: cfg,
			Logger: log,
		}

		var opt option.ClientOption

		// Check first for json format
		if serviceAccountRaw := os.Getenv("FIREBASE_SERVICE_ACCOUNT"); serviceAccountRaw != "" {
			// Try to parse as JSON directly
			if strings.HasPrefix(strings.TrimSpace(serviceAccountRaw), "{") {
				log.Info("Using Firebase credentials from environment variable parsed as JSON")
				opt = option.WithCredentialsJSON([]byte(serviceAccountRaw))
			} else {
				// Assume it's a file path
				if _, err := os.Stat(cfg.GetFirebaseCredentialsPath()); err == nil {
					log.Info("Using Firebase credentials from file", "path", cfg.GetFirebaseCredentialsPath())
					opt = option.WithCredentialsFile(cfg.GetFirebaseCredentialsPath())
				} else {
					log.Error("Firebase credentials file not found", "path", cfg.GetFirebaseCredentialsPath())
				}
			}
		} else {
			log.Warn("No Firebase credentials provided, attempting to use default credentials")
		}

		config := &firebase.Config{
			ProjectID: cfg.FirebaseProjectID,
		}

		// Firebase app initialization
		var app *firebase.App
		if opt != nil {
			app, err = firebase.NewApp(context.Background(), config, opt)
		} else {
			app, err = firebase.NewApp(context.Background(), config)
		}

		if err != nil {
			log.Error("Failed to initialize Firebase app", "error", err)
			return
		}

		instance.App = app

		// Initialize Auth
		instance.Auth, err = app.Auth(context.Background())
		if err != nil {
			log.Error("Failed to initialize Firebase Auth", "error", err)
		}

		// Initialize Firestore
		instance.Firestore, err = app.Firestore(context.Background())
		if err != nil {
			log.Error("Failed to initialize Firebase Firestore", "error", err)
		}

		// Initialize Storage
		instance.Storage, err = app.Storage(context.Background())
		if err != nil {
			log.Error("Failed to initialize Firebase Storage", "error", err)
		}

		log.Info("Firebase client initialized successfully", "projectID", cfg.FirebaseProjectID)
	})

	if instance.App == nil {
		return nil, err
	}

	return instance, nil
}

// GetClient returns the Firebase client instance
func GetClient() *Client {
	if instance == nil {
		panic("Firebase client not initialized")
	}
	return instance
}

// Close closes all Firebase connections
func (c *Client) Close() {
	if c.Firestore != nil {
		if err := c.Firestore.Close(); err != nil {
			c.Logger.Error("Failed to close Firestore client", "error", err)
		}
	}
}
