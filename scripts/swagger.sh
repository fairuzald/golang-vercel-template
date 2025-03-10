#!/bin/bash

# Show commands as they're executed
set -x

# Get go environment info
echo "Go environment:"
go env GOPATH
go env GOBIN

# Set GOBIN explicitly if not set
if [ -z "$(go env GOBIN)" ]; then
    export GOBIN=$(go env GOPATH)/bin
    echo "Set GOBIN to $GOBIN"
fi

# Add GOBIN to PATH if not already there
if [[ ":$PATH:" != *":$GOBIN:"* ]]; then
    export PATH="$GOBIN:$PATH"
    echo "Added $GOBIN to PATH"
fi

# Install swag with verbose output
echo "Installing swag..."
go install -v github.com/swaggo/swag/cmd/swag@latest

# Check if swag is available now
which swag || echo "swag not found in PATH: $PATH"

# Try with full path if not in PATH
SWAG_BIN="$GOBIN/swag"
if [ ! -f "$SWAG_BIN" ]; then
    echo "swag binary not found at $SWAG_BIN"
    # Try to find swag anywhere in GOPATH
    SWAG_BIN=$(find $(go env GOPATH) -name swag -type f | head -n 1)
    if [ -z "$SWAG_BIN" ]; then
        echo "Could not find swag binary anywhere in GOPATH"
        exit 1
    fi
    echo "Found swag at $SWAG_BIN"
fi

# Create docs directory if it doesn't exist
mkdir -p docs/swagger

# Fix the main.go location to be the actual main file with Swagger annotations
echo "Generating Swagger documentation..."
"$SWAG_BIN" init -g cmd/api/main.go --parseInternal --parseDependency -o docs/swagger

echo "Swagger documentation generated successfully!"
