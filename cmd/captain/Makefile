DIST := $(CURDIR)/_dist
BUILD := $(CURDIR)/_build
VERSION := ""
COMMIT := ""
LDFLAGS := "-X main.version=${VERSION} -X main.commit=${COMMIT}"
REGISTRY := softleader
BINARY = captain

.PHONY: build
build:
	mkdir -p $(BUILD)
	go build -o $(BUILD)/$(BINARY) ./cmd/$(BINARY)

.PHONY: dist
dist:
ifeq ($(strip $(VERSION)),)
	$(error VERSION is not set)
endif
ifeq ($(strip $(COMMIT)),)
	$(error COMMIT is not set)
endif
	mkdir -p $(DIST)
	cp ./cmd/$(BINARY)/Dockerfile $(DIST)
	cp ./build/icp-linux-amd64.tgz $(DIST)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(DIST)/$(BINARY) -ldflags $(LDFLAGS) -a -tags netgo ./cmd/$(BINARY)
	docker build -t $(REGISTRY)/$(BINARY) $(DIST)

.PHONY: ship
ship:
	 docker push $(REGISTRY)/$(BINARY)