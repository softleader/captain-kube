#!/usr/bin/env bash

docker pull hub.softleader.com.tw/captain-kube
docker-compose down
docker-compose up -d

exit 0
