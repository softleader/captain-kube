#!/bin/bash

make
docker build -t hub.softleader.com.tw/captain-kube .
docker push hub.softleader.com.tw/captain-kube