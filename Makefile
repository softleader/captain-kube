GOOS?=linux
GOARCH?=amd64
BINARY?=build
APP?=main

install:
	GOOS=${GOOS} GOARCH=${GOARCH} go build -o ${BINARY}/${APP} .