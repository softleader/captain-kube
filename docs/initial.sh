#!/bin/bash

cp ${PLAYBOOKS}/hosts /data
cp ${CAPTAIN_KUBE}/docker-compose.yml /data
cp /upgrade.sh /data

exit 0
