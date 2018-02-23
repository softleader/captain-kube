 #!/bin/bash

GOOS="$1"
GOARCH="$2"
APP=ck

if [[ -z $GOOS ]]; then
    GOOS=linux && echo "You did not determine GOOS, use \'$GOOS\' as default"
fi

if [[ -z $GOARCH ]]; then
    GOARCH=amd64 && echo "You did not determine GOARCH, use \'$GOARCH\' as default"
fi

if [[ $GOOS == "macos" ]]; then
    GOOS=darwin
fi

go get github.com/softleader/captain-kube && \
make -C $GOPATH/src/github.com/softleader/captain-kube GOOS=$GOOS $GOARCH=GOARCH BINARY=$(pwd)

exit 0