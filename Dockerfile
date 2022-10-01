## Build
FROM golang:1.19 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go build -o ./httpm

## Deploy
FROM gcr.io/distroless/base-debian11

WORKDIR /app

COPY --from=build /app/httpm .
COPY ./config.json .

ENTRYPOINT ["./httpm"]