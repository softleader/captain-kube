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

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: golint ## Run test
	go test ./... -v

.PHONY: gofmt
gofmt: ## Format source code
	gofmt -s -w .

.PHONY: golint
golint: gofmt ## Run linter for source code
ifndef HAS_GOLINT
	go get -u golang.org/x/lint/golint
endif
	golint -set_exit_status ./cmd/...
	golint -set_exit_status ./pkg/...

.PHONY: protoc
protoc: ## Parse proto files and generate output
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

.PHONY: build-caplet
build-caplet: ## Build Caplet binary
	make build --file ./cmd/$(CAPLET)/Makefile VERSION=$(VERSION) COMMIT=$(COMMIT) REGISTRY=$(REGISTRY)

.PHONY: build-captain
build-captain: ## Build Captain binary
	make build --file ./cmd/$(CAPLET)/Makefile VERSION=$(VERSION) COMMIT=$(COMMIT) REGISTRY=$(REGISTRY)

.PHONY: build-capui
build-capui: ## Build CapUI binary
	make build --file ./cmd/$(CAPUI)/Makefile VERSION=$(VERSION) COMMIT=$(COMMIT) REGISTRY=$(REGISTRY)

.PHONY: build-capctl
build-capctl: ## Build Capctl binary
	make build --file ./cmd/$(CAPCTL)/Makefile VERSION=$(VERSION) COMMIT=$(COMMIT)

.PHONY: dist-caplet
dist-caplet: ## Package Caplet to docker image
	make dist --file ./cmd/$(CAPLET)/Makefile VERSION=$(VERSION) COMMIT=$(COMMIT) REGISTRY=$(REGISTRY)

.PHONY: dist-captain
dist-captain: ## Package Captain to docker image
	make dist --file ./cmd/$(CAPTAIN)/Makefile VERSION=$(VERSION) COMMIT=$(COMMIT) REGISTRY=$(REGISTRY)

.PHONY: dist-capui
dist-capui: ## Package CapUI to docker image
	make dist --file ./cmd/$(CAPUI)/Makefile VERSION=$(VERSION) COMMIT=$(COMMIT) REGISTRY=$(REGISTRY)

.PHONY: dist-capctl
dist-capctl: ## Package Capctl to tarball
	make dist --file ./cmd/$(CAPCTL)/Makefile VERSION=$(VERSION) COMMIT=$(COMMIT)

.PHONY: ship-captain
ship-captain: ## Ship Captain docker image
	make ship --file ./cmd/$(CAPTAIN)/Makefile VERSION=$(VERSION) COMMIT=$(COMMIT) REGISTRY=$(REGISTRY)

.PHONY: ship-capui
ship-capui: ## Ship CapUI docker image
	make ship --file ./cmd/$(CAPUI)/Makefile VERSION=$(VERSION) COMMIT=$(COMMIT) REGISTRY=$(REGISTRY)

.PHONY: dist-caplet
ship-caplet: ## Ship Caplet docker image
	make ship --file ./cmd/$(CAPLET)/Makefile VERSION=$(VERSION) COMMIT=$(COMMIT)

.PHONE: release
release: release-caplet release-captain release-capui dist-capctl ## Dist and Ship all docker image and also package Capctl

.PHONE: release-caplet
release-caplet: dist-caplet ship-caplet ## Dist and Ship Caplet docker image

.PHONE: release-captain
release-captain: dist-caplet ship-captain  ## Dist and Ship Captain docker image

.PHONE: release-capui
release-capui: dist-caplet ship-capui ## Dist and Ship CapUI docker image

.PHONY: bootstrap
bootstrap: ## Fetch required go modules
	go mod download

.PHONY: clean
clean: ## Clean temporary folders
	rm -rf _*
