#!/bin/bash

GOOS="GOOS=linux"
GOARCH="GOARCH=amd64"
OUTPUT="ckube"

for arg
do
    if [ ${arg:0:4} = "GOOS" ]
    then
        GOOS=$arg
    elif [ ${arg:0:6} = "GOARCH" ]
    then
        GOARCH=$arg
    fi
done

if [ $GOOS = "GOOS=macos" ]
then
    GOOS="GOOS=darwin"
fi

go get github.com/softleader/captain-kube && \
make -C $GOPATH/src/github.com/softleader/captain-kube $GOOS $GOARCH BINARY=$(pwd) OUTPUT=$OUTPUT

exit 0
