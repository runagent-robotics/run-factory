FROM golang:1.22-alpine AS build

WORKDIR /src

COPY go.mod ./
COPY cmd ./cmd

RUN go build -o /out/runfactory-core ./cmd/server

FROM alpine:3.20

WORKDIR /app

COPY --from=build /out/runfactory-core /usr/local/bin/runfactory-core

EXPOSE 8080

ENTRYPOINT ["runfactory-core"]
