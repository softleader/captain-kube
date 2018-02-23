GOOS?=linux
GOARCH?=amd64
BINARY?=build
APP?=main

all:
    GOOS=darwin GOARCH=${GOARCH} go build -o ${BINARY}/${APP} .