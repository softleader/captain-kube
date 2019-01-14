FROM docker:18

RUN apk update && \
    apk add ca-certificates curl openssl && \
    rm -rf /var/cache/apk/* && \
    curl https://raw.githubusercontent.com/kubernetes/helm/master/scripts/get | sh && \
    helm init -c

WORKDIR /app
COPY docker-compose.yml /
COPY static/ /app/static/
COPY templates/ /app/templates/
COPY capui /app

ENV GIN_MODE=release

ENTRYPOINT ["./capui"]
