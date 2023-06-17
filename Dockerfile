FROM golang:alpine AS builder

WORKDIR /opt/src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build ./cmd/prometheus-multi-tenant-proxy-server

FROM alpine:latest

LABEL maintainer="Saeid Bostandoust <ssbostan@yahoo.com>"

EXPOSE 9999

WORKDIR /opt/app

COPY --from=builder /opt/src/prometheus-multi-tenant-proxy-server .
COPY --from=builder /opt/src/examples/namespace.yaml /opt/config/config.yaml

RUN adduser -DH rootless

USER 1000

ENTRYPOINT ["./prometheus-multi-tenant-proxy-server"]

CMD ["--config", "/opt/config/config.yaml"]
