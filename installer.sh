 #!/bin/bash

go get github.com/softleader/captain-kube && \
make -C $GOPATH/src/github.com/softleader/captain-kube "$1" BINARY=$(pwd)

exit 0