# 🚀 Golang Backend

An API service for a skill enhancement platform built with Go.

## 📋 Table of Contents

- [🚀 Golang Backend](#-golang-backend)
  - [📋 Table of Contents](#-table-of-contents)
  - [🔍 Overview](#-overview)
  - [✨ Features](#-features)
  - [🏗️ Project Structure](#️-project-structure)
  - [📝 Prerequisites](#-prerequisites)
  - [🚀 Getting Started](#-getting-started)
    - [Local Development](#local-development)
    - [Docker Development](#docker-development)
    - [Vercel Development](#vercel-development)
  - [🚢 Deployment Options](#-deployment-options)
    - [Docker Deployment](#docker-deployment)
    - [Vercel Deployment](#vercel-deployment)
  - [📚 API Documentation](#-api-documentation)
  - [⚙️ Configuration](#️-configuration)
  - [🔥 Firebase Integration](#-firebase-integration)
  - [🔄 Development Workflow](#-development-workflow)
  - [🔄 CI/CD Pipeline](#-cicd-pipeline)
  - [📝 Commands Reference](#-commands-reference)
  - [🔧 Troubleshooting](#-troubleshooting)
    - [Firebase Connection Issues](#firebase-connection-issues)
    - [Vercel Deployment Issues](#vercel-deployment-issues)
    - [Docker Issues](#docker-issues)
    - [API Issues](#api-issues)
  - [📄 License](#-license)

## 🔍 Overview

Golang Backend is a Go-based API service that provides the backend functionality for a skill enhancement platform. It uses Firebase for authentication, database, and storage, and can be deployed to Vercel for serverless operation.

## ✨ Features

- Clean architecture design
- Firebase authentication and data storage
- Health check monitoring
- API rate limiting
- CORS configuration
- Structured logging
- Swagger API documentation
- Docker support
- CI/CD with GitHub Actions
- Vercel deployment ready

## 🏗️ Project Structure

The project follows a domain-driven design with a clean architecture approach:

```
.
├── api/                  # API entry points and route definitions
├── app/                  # Application core
│   ├── core/             # Domain entities and interfaces
│   └── module/           # Feature modules
├── cmd/                  # Application entry points
├── configs/              # Configuration management
├── credentials/          # Credentials and secrets
├── docs/                 # Documentation, including Swagger
├── infrastructure/       # External services and tools
├── pkg/                  # Common packages and utilities
└── scripts/              # Utility scripts
```

## 📝 Prerequisites

- Go 1.21 or higher
- Docker (optional, for containerized development)
- Firebase project and credentials
- Vercel account (optional, for serverless deployment)
- Node.js and npm (optional, for Vercel CLI)

## 🚀 Getting Started

### Local Development

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/golang-vercel-template.git
   cd golang-vercel-template
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

3. Set up your Firebase credentials:

   - Create a Firebase project in Firebase Console
   - Generate a service account key and save it to `credentials/firebase-service-account.json`
   - Or set the `FIREBASE_SERVICE_ACCOUNT` environment variable with the service account JSON content

4. Create a `.env` file in the project root with the following variables (customize as needed):

   ```
   APP_NAME=golang-vercel-template
   APP_ENV=development
   APP_PORT=8080
   APP_DEBUG=true
   FIREBASE_PROJECT_ID=your-firebase-project-id
   ```

5. Run the application:

   ```bash
   go run cmd/api/main.go
   ```

   Or use hot reload with Air:

   ```bash
   make dev
   ```

### Docker Development

1. Make sure Docker and Docker Compose are installed on your system.

2. Set up your environment variables:

   - Copy `.env.example` to `.env` and fill in your Firebase credentials
   - Or provide environment variables directly in the `docker-compose.yml` file

3. Build and start the Docker container:

   ```bash
   make docker-up
   # or
   docker-compose up -d
   ```

4. Access the application at http://localhost:8080

5. View logs in real-time:

   ```bash
   make docker-logs
   # or
   docker-compose logs -f
   ```

6. Execute commands inside the container:

   ```bash
   docker-compose exec app sh
   # Then you can run commands inside the container
   go test ./...
   ```

7. Stop the container:

   ```bash
   make docker-down
   # or
   docker-compose down
   ```

8. Restart the container (after making changes):
   ```bash
   make docker-restart
   # or
   docker-compose restart
   ```

### Vercel Development

1. Install Vercel CLI globally:

   ```bash
   npm install -g vercel
   ```

2. Login to Vercel:

   ```bash
   vercel login
   ```

3. Link the project to your Vercel account:

   ```bash
   vercel link
   ```

4. Run the local development server:

   ```bash
   make vercel-dev
   # or
   vercel dev
   ```

5. Access the application at the URL provided by Vercel CLI (typically http://localhost:3000)

## 🚢 Deployment Options

### Docker Deployment

1. Build a production Docker image:

   ```bash
   docker build -t golang-vercel-template:prod --target production .
   ```

2. Run the container:

   ```bash
   docker run -d -p 8080:8080 \
     -e APP_ENV=production \
     -e FIREBASE_PROJECT_ID=your-project-id \
     -e FIREBASE_SERVICE_ACCOUNT='{"type":"service_account",...}' \
     --name golang-vercel-template golang-vercel-template:prod
   ```

3. For Kubernetes deployment, you can use the provided Kubernetes manifests:

   ```bash
   kubectl apply -f kubernetes/
   ```

### Vercel Deployment

1. Install the Vercel CLI if not already installed:

   ```bash
   npm install -g vercel
   ```

2. Run the flattening script to prepare for Vercel:

   ```bash
   ./scripts/flatten.sh
   ```

3. Deploy to Vercel:

   ```bash
   cd out
   vercel
   ```

   Or deploy directly using the provided command:

   ```bash
   make vercel
   ```

4. For production deployment:

   ```bash
   vercel --prod
   ```

5. Set up environment variables in the Vercel dashboard:
   - Go to your project in the Vercel dashboard
   - Navigate to Settings > Environment Variables
   - Add all required environment variables (APP_ENV, FIREBASE_PROJECT_ID, etc.)

## 📚 API Documentation

Swagger documentation is available in development mode:

1. Generate the Swagger documentation:

   ```bash
   make docs
   ```

2. Access the Swagger UI at:
   ```
   http://localhost:8080/swagger/index.html
   ```

## ⚙️ Configuration

Configuration is managed through environment variables and the `configs/config.go` file. You can set the following environment variables:

| Variable                 | Description                          | Default                                     |
| ------------------------ | ------------------------------------ | ------------------------------------------- |
| APP_NAME                 | Application name                     | golang-vercel-template                      |
| APP_ENV                  | Environment (development/production) | development                                 |
| APP_PORT                 | HTTP server port                     | 8080                                        |
| APP_DEBUG                | Enable debug logging                 | true                                        |
| APP_SECRET               | Secret key for encryption/JWT        | your-secret-key-at-least-32-chars-long      |
| FIREBASE_PROJECT_ID      | Firebase project ID                  | -                                           |
| FIREBASE_SERVICE_ACCOUNT | Firebase credentials JSON            | ./credentials/firebase-service-account.json |
| RATE_LIMIT_REQUESTS      | API rate limit requests count        | 100                                         |
| RATE_LIMIT_DURATION      | API rate limit duration              | 1m                                          |
| CORS_ALLOWED_ORIGINS     | CORS allowed origins                 | \*                                          |
| CORS_ALLOWED_METHODS     | CORS allowed methods                 | GET,POST,PUT,DELETE,OPTIONS                 |
| CORS_ALLOWED_HEADERS     | CORS allowed headers                 | Authorization,Content-Type,X-Requested-With |
| CORS_EXPOSED_HEADERS     | CORS exposed headers                 | Content-Length                              |
| CORS_MAX_AGE             | CORS preflight max age               | 12h                                         |
| AUTH_TOKEN_EXPIRY        | JWT token expiry                     | 24h                                         |

## 🔥 Firebase Integration

This project uses Firebase for:

- Authentication
- Firestore database
- Cloud Storage

To set up Firebase:

1. Create a Firebase project at https://console.firebase.google.com/
2. Generate a new service account key from Project Settings > Service Accounts
3. Save the key to `credentials/firebase-service-account.json` or use the environment variable
4. Enable the services you need (Authentication, Firestore, Storage) in the Firebase console
5. Set the `FIREBASE_PROJECT_ID` environment variable to your project ID

## 🔄 Development Workflow

1. Create a feature branch from `dev`:

   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes and commit them.

3. Push your branch and create a PR to the `dev` branch.

4. The CI/CD pipeline will automatically deploy your changes to a preview URL.

5. Once approved and merged to `dev`, it will be deployed to the development environment.

6. Releases to production are made by merging `dev` into `main`.

## 🔄 CI/CD Pipeline

The project includes GitHub Actions workflows for CI/CD:

- `dev.yml`: Runs on pull requests to the `dev` branch and pushes to `dev`

  - Deploys to a Vercel preview environment
  - Adds deployment comments to PRs and issues

- `prod.yml`: Runs on pushes to the `main` branch
  - Deploys to the Vercel production environment
  - Creates deployment history issues

## 📝 Commands Reference

The project includes a Makefile with useful commands:

```bash
# Build the application
make build

# Run the application locally
make run

# Run with hot reload
make dev

# Generate Swagger documentation
make docs

# Start Docker for development
make docker-up

# View Docker logs
make docker-logs

# Stop Docker containers
make docker-down

# Restart Docker containers
make docker-restart

# Clean Docker environment
make docker-clean

# Deploy to Vercel
make vercel

# Run Vercel dev environment
make vercel-dev

# Show all available commands
make help
```

## 🔧 Troubleshooting

### Firebase Connection Issues

- Verify your Firebase service account key is correct
- Check that your `FIREBASE_PROJECT_ID` matches the project ID in your service account key
- Ensure the Firebase project has Firestore and Authentication enabled
- Check logs for Firebase connection errors: `make docker-logs`

### Vercel Deployment Issues

- Make sure you've run the flattening script: `./scripts/flatten.sh`
- Verify the `vercel.json` configuration points to the correct entry point
- Check that all required environment variables are set in Vercel
- Inspect the deployment logs: `vercel logs`

### Docker Issues

- Check if Docker and Docker Compose are installed and running:
  ```bash
  docker --version
  docker-compose --version
  ```
- Verify that no other services are using port 8080:
  ```bash
  lsof -i :8080
  ```
- Check container logs for errors:
  ```bash
  docker-compose logs
  ```
- Try rebuilding the container from scratch:
  ```bash
  make docker-clean
  make docker-up
  ```

### API Issues

- Check the logs for detailed error messages
- Use the `/api/health` endpoint to verify service status
- Ensure CORS is configured correctly for your frontend
- Verify that your Firebase credentials have the correct permissions

## 📄 License

This project is licensed under the MIT License.
