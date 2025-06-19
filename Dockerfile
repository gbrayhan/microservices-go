FROM golang:1.24-alpine AS builder

WORKDIR /srv/go-app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -installsuffix cgo -o microservice .

FROM gcr.io/distroless/static:nonroot

WORKDIR /srv/go-app
COPY --from=builder /srv/go-app/microservice .

USER nonroot:nonroot

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s \
  CMD ["/srv/go-app/microservice","health"]

CMD ["./microservice"]
