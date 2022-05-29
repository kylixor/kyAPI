FROM golang:alpine

WORKDIR /go/src/app

RUN go install github.com/cespare/reflex@latest
LABEL org.opencontainers.image.description "https://github.com/kylixor/kyAPI-dev"