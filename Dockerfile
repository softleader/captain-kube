FROM golang:alpine
MAINTAINER softleader.com.tw

RUN apk update && \
	apk --no-cache add bash make git && \
	rm -rf /var/cache/apk/*

COPY installer.sh /installer.sh

WORKDIR /data

ENTRYPOINT ["/bin/sh", "/installer.sh"]