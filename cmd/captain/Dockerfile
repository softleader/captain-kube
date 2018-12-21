FROM alpine:3.8

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /app
ADD captain /app/captain

ENV CAPTAIN_PORT=80
ENTRYPOINT ./captain
