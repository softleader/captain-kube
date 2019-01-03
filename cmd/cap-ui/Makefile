DIST := $(CURDIR)/_dist
BUILD := $(CURDIR)/_build
VERSION := ""
COMMIT := ""
LDFLAGS := "-X main.version=${VERSION} -X main.commit=${COMMIT}"
REGISTRY := softleader
BINARY = cap-ui

.PHONY: build
build:
	mkdir $(BUILD)
	go build -o $(BUILD)/$(BINARY) ./cmd/$(BINARY)

.PHONY: dist
dist:
ifeq ($(strip $(VERSION)),)
	$(error VERSION is not set)
endif
ifeq ($(strip $(COMMIT)),)
	$(error COMMIT is not set)
endif
	mkdir "$(DIST)"
	cp ./cmd/$(BINARY)/Dockerfile $(DIST)
	cp -r ./templates $(DIST)
	cp -r ./static $(DIST)
	cp ./cmd/$(BINARY)/docker-compose.yml $(DIST)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(DIST)/$(BINARY) -ldflags $(LDFLAGS) -a -tags netgo ./cmd/$(BINARY)
	docker build -t $(REGISTRY)/$(BINARY) $(DIST) && docker push $(REGISTRY)/$(BINARY)
