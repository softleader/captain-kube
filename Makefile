HAS_GLIDE := $(shell command -v glide;)
DIST := $(CURDIR)/_dist
BUILD := $(CURDIR)/_build
BINARY := caplet
REGISTRY := softleader

.PHONY: install
install: bootstrap test build
	mkdir -p $(SL_PLUGIN_DIR)
	cp $(BUILD)/$(BINARY) $(SL_PLUGIN_DIR)
	cp $(METADATA) $(SL_PLUGIN_DIR)

.PHONY: test
test:
	go test ./... -v

.PHONY: build
build: clean bootstrap
	mkdir -p $(BUILD)
	go build -o $(BUILD)/$(BINARY)

.PHONY: dist
dist:
	mkdir -p $(DIST)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(DIST)/$(BINARY) -a -tags netgo $(MAIN)
	docker build -t $(REGISTRY)/$(BINARY) .
	docker push $(REGISTRY)/$(BINARY)

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
