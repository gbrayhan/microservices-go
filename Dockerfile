FROM golang:1.14 AS builder
WORKDIR /srv
ADD . .
RUN go build -o microservice


FROM debian:buster
WORKDIR /srv
COPY --from=builder /srv/config.json .
COPY --from=builder /srv/views ./views/
COPY --from=builder /srv/archives ./archives/
COPY --from=builder /srv/public ./public/
COPY --from=builder /srv/microservice .

CMD ["./microservice"]