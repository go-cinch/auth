FROM golang:1.18 AS builder

COPY . /src
WORKDIR /src

RUN make build

FROM debian:stable-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates netbase && \
    rm -rf /var/lib/apt/lists/ && \
    apt-get autoremove -y && \
    apt-get autoclean -y

COPY --from=builder /src/bin /app

WORKDIR /app

EXPOSE 6060
EXPOSE 6160
COPY configs /data/conf

CMD ["sh", "-c", "./auth -c /data/conf"]
# enable k8s config map
#CMD ["sh", "-c", "./auth -c /data/conf -n cinch -l auth"]
