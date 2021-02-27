FROM golang:1.14 AS builder
WORKDIR /srv/go-app
ADD . .
RUN go build -o microservice


FROM debian:buster
WORKDIR /srv/go-app
COPY --from=builder /srv/go-app/config.json .
COPY --from=builder /srv/go-app/views ./views/
COPY --from=builder /srv/go-app/archives ./archives/
COPY --from=builder /srv/go-app/public ./public/
COPY --from=builder /srv/go-app/microservice .

CMD ["./microservice"]