FROM golang:1.19 AS builder
ADD . /app/cli
WORKDIR /app/cli
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /usr/local/bin/kafka-topic-creator .
