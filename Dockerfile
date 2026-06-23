# --- Build stage ---
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Cache dependencies first
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o mininexus ./cmd/api

# --- Runtime stage ---
FROM gcr.io/distroless/static-debian12

WORKDIR /app
COPY --from=builder /app/mininexus .

# Cloud Run injects PORT; default to 8080
ENV PORT=8080
EXPOSE 8080

USER nonroot:nonroot
ENTRYPOINT ["/app/mininexus"]
