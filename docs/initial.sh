#!/bin/bash

cp ${PLAYBOOKS}/hosts /data
cp /docker-compose.yml /data
cp /daemon.yaml /data
cp /upgrade.sh /data

exit 0
