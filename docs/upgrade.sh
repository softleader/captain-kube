#!/usr/bin/env bash

docker pull softleader/captain-kube | grep -q 'Image is up to date'
# docker pull hub.softleader.com.tw/captain-kube | grep -q 'Image is up to date'

if [ $? -eq 0 ]; then
    echo "Image already up to date"
else
    docker-compose down
    docker-compose up -d
fi

exit 0
