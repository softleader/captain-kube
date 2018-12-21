HAS_GLIDE := $(shell command -v glide;)
DIST := $(CURDIR)/_dist
BUILD := $(CURDIR)/_build
REGISTRY := softleader

CAPTAIN := captain
CAPLET := caplet
UI := cap-ui
CAPCTL = capctl

.PHONY: test
test:
	go test ./... -v

build: protoc
protoc:
	protoc -I api/protobuf-spec/ --go_out=plugins=grpc:pkg/proto api/protobuf-spec/*.proto

.PHONY: build
build: clean bootstrap build-caplet build-captain build-ui build-capctl

.PHONY: build-caplet
build-caplet:
	make build --file ./cmd/$(CAPLET)/Makefile REGISTRY=$(REGISTRY)

.PHONY: build-captain
build-captain:
	make build --file ./cmd/$(CAPLET)/Makefile REGISTRY=$(REGISTRY)

.PHONY: build-ui
build-ui:
	make build --file ./cmd/$(UI)/Makefile REGISTRY=$(REGISTRY)

.PHONY: build-capctl
build-capctl:
	make build --file ./cmd/$(CAPCTL)/Makefile REGISTRY=$(REGISTRY)

.PHONY: dist
dist: dist-caplet dist-captain dist-ui dist-calctl

.PHONY: dist-caplet
dist-caplet:
	make dist --file ./cmd/$(CAPLET)/Makefile REGISTRY=$(REGISTRY)

.PHONY: dist-captain
dist-captain:
	make dist --file ./cmd/$(CAPTAIN)/Makefile REGISTRY=$(REGISTRY)

.PHONY: dist-ui
dist-ui:
	make dist --file ./cmd/$(UI)/Makefile REGISTRY=$(REGISTRY)

.PHONY: dist-calctl
dist-calctl:
	make dist --file ./cmd/$(CAPCTL)/Makefile REGISTRY=$(REGISTRY)

.PHONY: bootstrap
bootstrap:
ifndef HAS_GLIDE
	go get -u github.com/Masterminds/glide
endif
ifeq (,$(wildcard ./glide.yaml))
	glide init
endif
	glide install --strip-vendor

.PHONY: clean
clean:
	rm -rf _*