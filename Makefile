HAS_GLIDE := $(shell command -v glide;)
DIST := $(CURDIR)/_dist
BUILD := $(CURDIR)/_build
REGISTRY := softleader
CAPTAIN := captain
CAPLET := caplet

.PHONY: test
test:
	go test ./... -v

.PHONY: build
build: clean bootstrap
	mkdir -p $(BUILD)
	go build -o $(BUILD)/$(BINARY) ./cmd/$(CAPTAIN)
	go build -o $(BUILD)/$(BINARY) ./cmd/$(CAPLET)

.PHONY: dist
dist:
	mkdir -p $(DIST)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(DIST)/$(BINARY) -a -tags netgo $(CAPTAIN)
	docker build -t $(REGISTRY)/$(CAPTAIN) . && docker push $(REGISTRY)/$(CAPTAIN)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(DIST)/$(BINARY) -a -tags netgo $(CAPLET)
	docker build -t $(REGISTRY)/$(CAPLET) . && docker push $(REGISTRY)/$(CAPLET)

.PHONY: bootstrap
bootstrap:
ifndef HAS_GLIDE
	go get -u github.com/Masterminds/glide
endif
ifeq (,$(wildcard ./glide.yaml))
	glide init --non-interactive
endif
	glide install --strip-vendor	

.PHONY: clean
clean:
	rm -rf _*