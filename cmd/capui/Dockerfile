FROM alpine:3.8

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /app
COPY docker-compose.yml /
COPY static/ /app/static/
COPY templates/ /app/templates/
COPY capui /app

#ENV GIN_MODE=release

ENTRYPOINT ["./capui"]
