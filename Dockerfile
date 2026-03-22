# syntax=docker/dockerfile:1
ARG GO_VERSION=1.23
# Note: User had 1.25.5-alpine, but 1.23 is more standard for current projects. 
# I'll stick to 1.25.5-alpine if they really want it, but I'll use a variable.
ARG IMAGE_TAG=1.25.5-alpine

################################################################################
# Base stage for shared dependencies
FROM golang:${IMAGE_TAG} AS base
RUN apk add --no-cache git make
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

################################################################################
# Development stage
FROM base AS dev
RUN go install github.com/air-verse/air@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest
COPY . .
# Docs generation for development
RUN swag init -g main.go --parseDependency --parseInternal \
    --dir ./cmd/api,./internal/features/auth,./internal/features/users,./internal/features/appointments,./internal/features/excuseslips,./internal/features/students \
    --output ./docs/internal --instanceName internal
RUN swag init -g main.go --parseDependency --parseInternal \
    --dir ./cmd/api,./internal/features/students/external \
    --output ./docs/external --instanceName external
CMD ["air"]

################################################################################
# Builder stage for production
FROM base AS builder
RUN go install github.com/swaggo/swag/cmd/swag@latest
COPY . .
RUN swag init -g main.go --parseDependency --parseInternal \
    --dir ./cmd/api,./internal/features/auth,./internal/features/users,./internal/features/appointments,./internal/features/excuseslips,./internal/features/students \
    --output ./docs/internal --instanceName internal
RUN swag init -g main.go --parseDependency --parseInternal \
    --dir ./cmd/api,./internal/features/students/external \
    --output ./docs/external --instanceName external
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

################################################################################
# Production stage
FROM alpine:latest AS prod
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/.env.example ./.env
# Note: You should ideally provide a real .env or use environment variables in compose
EXPOSE 8080
CMD ["./main"]
