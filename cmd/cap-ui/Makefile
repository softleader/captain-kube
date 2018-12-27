DIST := $(CURDIR)/_dist
BUILD := $(CURDIR)/_build
REGISTRY := softleader
BINARY = cap-ui

.PHONY: build
build:
	mkdir $(BUILD)
	go build -o $(BUILD)/$(BINARY) ./cmd/$(BINARY)

.PHONY: dist
dist:
	mkdir "$(DIST)"
	cp ./cmd/$(BINARY)/Dockerfile $(DIST)
	cp -r ./templates $(DIST)
	cp -r ./static $(DIST)
	cp ./cmd/$(BINARY)/docker-compose.yml $(DIST)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(DIST)/$(BINARY) -a -tags netgo ./cmd/$(BINARY)
	docker build -t $(REGISTRY)/$(BINARY) $(DIST) && docker push $(REGISTRY)/$(BINARY)
