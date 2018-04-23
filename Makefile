GOOS?=linux
GOARCH?=amd64
BINARY?=dist
OUTPUT?=main

install:
	GOOS=${GOOS} GOARCH=${GOARCH} go build -o ${BINARY}/${OUTPUT} .