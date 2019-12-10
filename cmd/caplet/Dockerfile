FROM docker:18

ARG GRPC_HEALTH_PROBE_VERSION=v0.3.1

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/* && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

WORKDIR /app
ADD caplet /app/caplet

ENV GIN_MODE=release
ENV CAPLET_PORT=50051

ENTRYPOINT ["./caplet"]