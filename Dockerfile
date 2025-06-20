FROM golang:1.24-alpine AS builder

WORKDIR /srv/go-app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -installsuffix cgo -o microservice .

FROM alpine:3.20

WORKDIR /srv/go-app
COPY --from=builder /srv/go-app/microservice .

# Install curl for healthcheck
RUN apk add --no-cache curl

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

USER appuser:appgroup

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/v1/user/ || exit 1

CMD ["./microservice"]
