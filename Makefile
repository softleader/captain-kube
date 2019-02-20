HAS_GOLINT := $(shell command -v golint;)
DIST := $(CURDIR)/_dist
BUILD := $(CURDIR)/_build
REGISTRY := softleader
VERSION := ""
COMMIT := ""
CAPTAIN := captain
CAPLET := caplet
CAPUI := capui
CAPCTL = capctl
protobuf = api/protobuf-spec/softleader/captainkube/v2
proto_dst = pkg/proto


.PHONY: test
test: golint
	go test ./... -v

.PHONY: gofmt
gofmt:
	gofmt -s -w .

.PHONY: golint
golint: gofmt
ifndef HAS_GOLINT
	go get -u golang.org/x/lint/golint
endif
	golint -set_exit_status ./cmd/...
	golint -set_exit_status ./pkg/...

.PHONY: protoc
protoc:
	protoc -I $(protobuf) --go_out=plugins=grpc:$(proto_dst) $(protobuf)/caplet.proto
	protoc -I $(protobuf) --go_out=plugins=grpc:$(proto_dst) $(protobuf)/captain.proto
	protoc -I $(protobuf) --go_out=plugins=grpc:$(proto_dst) $(protobuf)/chart.proto
	protoc -I $(protobuf) --go_out=plugins=grpc:$(proto_dst) $(protobuf)/image.proto
	protoc -I $(protobuf) --go_out=plugins=grpc:$(proto_dst) $(protobuf)/msg.proto
	protoc -I $(protobuf) --go_out=plugins=grpc:$(proto_dst) $(protobuf)/prune.proto
	protoc -I $(protobuf) --go_out=plugins=grpc:$(proto_dst) $(protobuf)/version.proto
	protoc -I $(protobuf) --go_out=plugins=grpc:$(proto_dst) $(protobuf)/console_url.proto
	protoc -I $(protobuf) --go_out=plugins=grpc:$(proto_dst) $(protobuf)/rmi.proto
	protoc -I $(protobuf) --go_out=plugins=grpc:$(proto_dst) $(protobuf)/rmc.proto

.PHONY: build
build: clean bootstrap build-caplet build-captain build-capui build-capctl

.PHONY: build-caplet
build-caplet:
	make build --file ./cmd/$(CAPLET)/Makefile VERSION=$(VERSION) COMMIT=$(COMMIT) REGISTRY=$(REGISTRY)

.PHONY: build-captain
build-captain:
	make build --file ./cmd/$(CAPLET)/Makefile VERSION=$(VERSION) COMMIT=$(COMMIT) REGISTRY=$(REGISTRY)

.PHONY: build-capui
build-capui:
	make build --file ./cmd/$(CAPUI)/Makefile VERSION=$(VERSION) COMMIT=$(COMMIT) REGISTRY=$(REGISTRY)

.PHONY: build-capctl
build-capctl:
	make build --file ./cmd/$(CAPCTL)/Makefile VERSION=$(VERSION) COMMIT=$(COMMIT)

.PHONY: dist
dist: dist-caplet dist-captain dist-capui dist-capctl

.PHONY: dist-caplet
dist-caplet:
	make dist --file ./cmd/$(CAPLET)/Makefile VERSION=$(VERSION) COMMIT=$(COMMIT) REGISTRY=$(REGISTRY)

.PHONY: dist-captain
dist-captain:
	make dist --file ./cmd/$(CAPTAIN)/Makefile VERSION=$(VERSION) COMMIT=$(COMMIT) REGISTRY=$(REGISTRY)

.PHONY: dist-capui
dist-capui:
	make dist --file ./cmd/$(CAPUI)/Makefile VERSION=$(VERSION) COMMIT=$(COMMIT) REGISTRY=$(REGISTRY)

.PHONY: dist-capctl
dist-capctl:
	make dist --file ./cmd/$(CAPCTL)/Makefile VERSION=$(VERSION) COMMIT=$(COMMIT)

.PHONY: bootstrap
bootstrap:
	go mod download

.PHONY: clean
clean:
	rm -rf _*
