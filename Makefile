GOOS?=linux
GOARCH?=amd64
BINARY?=dist
OUTPUT?=main

.PHONY: all clean

all: test build publish clean

test:
	go test -cover -race ./...

build:
	GOOS=${GOOS} GOARCH=${GOARCH} go build -o ${BINARY}/${OUTPUT} .

publish:
	docker build -t hub.softleader.com.tw/captain-kube .
	docker push hub.softleader.com.tw/captain-kube
	docker tag hub.softleader.com.tw/captain-kube softleader/captain-kube
	docker push softleader/captain-kube

clean:
	rm -rf ${BINARY}/