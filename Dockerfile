FROM golang:1.13.1

RUN set -eux; \
    mkdir /build
WORKDIR /build

COPY go.mod go.sum ./
RUN set -eux; \
    go mod download

COPY fixtures fixtures
COPY *.go ./

RUN go test -v
