FROM golang:alpine as build

RUN apk add ca-certificates git gcc musl-dev mingw-w64-gcc

WORKDIR /opt

COPY go.mod go.sum ./
RUN  go mod download

COPY cmd/client      cmd/client
COPY internal/pow    internal/pow
COPY internal/client internal/client

RUN cd /opt/cmd/client && \
    go build -o /srv/client


FROM alpine:latest

COPY --from=build /srv /srv

WORKDIR /srv
CMD /srv/client
