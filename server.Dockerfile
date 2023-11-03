FROM golang:alpine as build

RUN apk add ca-certificates git gcc musl-dev mingw-w64-gcc

WORKDIR /opt

COPY go.mod go.sum ./
RUN  go mod download

COPY cmd/server      cmd/server
COPY internal/pow    internal/pow
COPY internal/server internal/server

RUN cd /opt/cmd/server && \
    go build -o /srv/server


FROM alpine:latest

COPY --from=build /srv /srv

WORKDIR /srv
CMD /srv/server
