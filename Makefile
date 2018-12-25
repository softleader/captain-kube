HAS_DEP := $(shell command -v dep;)
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
	protoc -I api/protobuf-spec/ --go_out=plugins=grpc:pkg/proto api/protobuf-spec/caplet.proto
	protoc -I api/protobuf-spec/ --go_out=plugins=grpc:pkg/proto api/protobuf-spec/captain.proto
	protoc -I api/protobuf-spec/ --go_out=plugins=grpc:pkg/proto api/protobuf-spec/chart.proto
	protoc -I api/protobuf-spec/ --go_out=plugins=grpc:pkg/proto api/protobuf-spec/image.proto
	protoc -I api/protobuf-spec/ --go_out=plugins=grpc:pkg/proto api/protobuf-spec/msg.proto
	protoc -I api/protobuf-spec/ --go_out=plugins=grpc:pkg/proto api/protobuf-spec/prune.proto
	protoc -I api/protobuf-spec/ --go_out=plugins=grpc:pkg/proto api/protobuf-spec/version.proto

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
ifndef HAS_DEP
	go get -u github.com/golang/dep/cmd/dep
endif
ifeq (,$(wildcard ./Gopkg.toml))
	dep init
endif
	dep ensure

.PHONY: clean
clean:
	rm -rf _*
