#!/bin/bash

cp ${PLAYBOOKS}/hosts /data
cp -R /nginx/* /data
cp /docker-compose.yml /data
cp /daemon.yaml /data
cp /upgrade.sh /data

exit 0
