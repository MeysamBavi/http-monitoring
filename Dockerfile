## Build
FROM golang:latest AS build

ENV GOPATH=/app
WORKDIR /app/src/http-montoring

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./internal/* ./internal/
COPY main.go ./

RUN go build -o ./httpm

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /app/src/http-montoring/httpm /
COPY ./config.json /

ENTRYPOINT ["/httpm"]