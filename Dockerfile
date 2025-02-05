FROM golang:1.20 AS builder
WORKDIR /srv/go-app
COPY . .
RUN go build -o microservice


FROM golang:1.20
WORKDIR /srv/go-app
#COPY --from=builder /srv/go-app/other-archives ./other-archives/
COPY --from=builder /srv/go-app/microservice .

CMD ["./microservice"]