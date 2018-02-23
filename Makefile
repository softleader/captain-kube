GOARCH?=amd64
BINARY?=build
APP?=ck

all: clean macos linux windows

macos:
	GOOS=darwin GOARCH=${GOARCH} go build -o ${BINARY}/${APP} .

linux:
	GOOS=linux GOARCH=${GOARCH} go build -o ${BINARY}/${APP} .

windows:
	GOOS=windows GOARCH=${GOARCH} go build -o ${BINARY}/${APP}.exe .

clean:
	rm -rf ${BINARY}

.PHONY: clean