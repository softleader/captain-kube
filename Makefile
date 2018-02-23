GOOS?=linux
GOARCH?=amd64
BINARY?=build
OUTPUT?=main

install:
	GOOS=${GOOS} GOARCH=${GOARCH} go build -o ${BINARY}/${OUTPUT} .