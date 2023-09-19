FROM golang:1.20 AS builder
WORKDIR /srv/go-app
COPY . .
RUN go build -o microservice


FROM golang:1.20
WORKDIR /srv/go-app
COPY --from=builder /srv/go-app/config.json .
COPY --from=builder /srv/go-app/archives ./archives/
COPY --from=builder /srv/go-app/microservice .

CMD ["./microservice"]