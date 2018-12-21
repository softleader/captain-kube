FROM alpine:3.8

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /app
ADD ui /app/ui
ADD static/* /app/static/
ADD templates/* /app/templates/
ADD templates/* /app/templates/
ADD configs/* /app/configs/

ENTRYPOINT ./ui