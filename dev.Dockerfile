# syntax=docker/dockerfile:1

FROM golang:1.19-alpine as BuildStage

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o main cmd/server/main.go

FROM alpine:latest

WORKDIR /

COPY --from=BuildStage app/configs/config.yml /configs/config.yml
COPY --from=BuildStage app/web /web
COPY --from=BuildStage app/assets /assets
COPY --from=BuildStage app/main /main

CMD ["./main"]