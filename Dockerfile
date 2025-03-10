FROM golang:1.21

RUN go install github.com/cosmtrek/air@v1.44.0

RUN go install github.com/swaggo/swag/cmd/swag@v1.16.2

RUN go install github.com/go-delve/delve/cmd/dlv@v1.21.0

RUN go install golang.org/x/tools/gopls@v0.14.2

RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY .air.toml ./

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]
