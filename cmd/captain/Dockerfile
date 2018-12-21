FROM alpine:3.8

RUN apk update && \
    apk add curl openssl ca-certificates && \
    rm -rf /var/cache/apk/* && \
	curl https://raw.githubusercontent.com/kubernetes/helm/master/scripts/get | sh && \
	helm init -c

WORKDIR /app
ADD captain /app/captain

ENTRYPOINT ["./captain"]
