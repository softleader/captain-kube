 #!/bin/bash

GOOS="$1"
GOARCH="$2"
APP=ck

if [[ -z $GOOS ]]; then
    GOOS=linux && echo "Use '$GOOS' as default GOOS"
fi

if [[ -z $GOARCH ]]; then
    GOARCH=amd64 && echo "Use '$GOARCH' as default GOARCH"
fi

if [[ $GOOS == "macos" ]]; then
    GOOS=darwin
fi

go get github.com/softleader/captain-kube && \
make -C $GOPATH/src/github.com/softleader/captain-kube GOOS=$GOOS $GOARCH=GOARCH APP=$APP BINARY=$(pwd)

exit 0