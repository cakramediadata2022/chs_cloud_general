# syntax=docker/dockerfile:1
# Initial stage: download modules
FROM golang:alpine as builder

RUN apk add build-base

ENV config=docker

WORKDIR /app

COPY ./ /app

RUN go mod download

RUN go mod tidy
# Intermediate stage: Build the binary
FROM golang:alpine as runner

COPY --from=builder ./app ./app

RUN go install -mod=mod github.com/githubnemo/CompileDaemon

RUN apk add build-base

WORKDIR /app
ENV config=docker

EXPOSE 9000

ENTRYPOINT ["CompileDaemon", "--build=go build cmd/server/main.go", "--command=./main"] 
 