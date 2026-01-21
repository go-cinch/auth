FROM golang:1.25-alpine3.23 AS builder

#ENV GOPROXY=https://goproxy.cn
# CGO_ENABLED=1, need check `ldd --version` is same as builder
ENV CGO_ENABLED=0

RUN apk add --no-cache git make bash

WORKDIR /src
# download first can use docker build cache if go.mod not change
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

COPY . .
RUN make build

FROM alpine:3.23

COPY --from=builder /src/bin /app

WORKDIR /app

COPY configs /data/conf

CMD ["sh", "-c", "./auth -c /data/conf"]
