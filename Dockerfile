FROM golang:1.19 AS builder
ADD . .
RUN go install .
