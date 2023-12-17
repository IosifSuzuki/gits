# syntax=docker/dockerfile:1

FROM golang:1.19 as BuildStage

RUN CGO_ENABLED=0 go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 go build -gcflags="all=-N -l" -o main cmd/server/main.go

FROM debian:buster

WORKDIR /

COPY --from=BuildStage /go/bin/dlv /
COPY --from=BuildStage app/configs/config.yml /configs/config.yml
COPY --from=BuildStage app/web /web
COPY --from=BuildStage app/assets /assets
COPY --from=BuildStage app/main /main

CMD ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "./main"]