FROM alpine:3.8

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /app
ADD caplet /app/caplet

ENV GIN_MODE=release
ENV CAPLET_PORT=50051

ENTRYPOINT ["./caplet"]
