DIST := $(CURDIR)/_dist
BUILD := $(CURDIR)/_build
REGISTRY := softleader
BINARY = ui

.PHONY: build
build:
	mkdir -p $(BUILD)
	go build -o $(BUILD)/$(BINARY) ./cmd/$(BINARY)

.PHONY: dist
dist:
	mkdir -p $(DIST)
	cp ./cmd/$(BINARY)/Dockerfile $(DIST)
	cp ./templates $(DIST)
	cp ./configs $(DIST)
	cp ./static $(DIST)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(DIST)/$(BINARY) -a -tags netgo ./cmd/$(BINARY)
	docker build -t $(REGISTRY)/$(BINARY) $(DIST) && docker push $(REGISTRY)/$(BINARY)