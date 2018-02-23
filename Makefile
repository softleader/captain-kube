GOARCH=amd64
BINARY=build
APP=captain-kube

all: clean macos linux windows

macos:
	GOOS=darwin GOARCH=${GOARCH} go build -o ${BINARY}/${APP}-macos-${GOARCH} .

linux:
	GOOS=linux GOARCH=${GOARCH} go build -o ${BINARY}/${APP}-linux-${GOARCH} .

windows:
	GOOS=windows GOARCH=${GOARCH} go build -o ${BINARY}/${APP}-windows-${GOARCH}.exe .

clean:
	rm -rf ${BINARY}

.PHONY: clean