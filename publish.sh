#!/bin/bash

make
docker build -t hub.softleader.com.tw/captain-kube .
docker push hub.softleader.com.tw/captain-kube
docker tag hub.softleader.com.tw/captain-kube softleader/captain-kube
docker push softleader/captain-kube