GOOS?=linux
GOARCH?=amd64
BINARY?=dist
OUTPUT?=main

all: test build

test:
	go test -cover -race ./...

build:
	GOOS=${GOOS} GOARCH=${GOARCH} go build -o ${BINARY}/${OUTPUT} .
