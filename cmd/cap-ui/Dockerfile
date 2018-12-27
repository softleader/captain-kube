FROM alpine:3.8

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /app
ADD docker-compose.yml /data/
ADD static/* /app/static/
ADD templates/* /app/templates/

ENTRYPOINT ["./ui"]
