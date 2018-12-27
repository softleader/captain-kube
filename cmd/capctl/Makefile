HAS_GLIDE := $(shell command -v glide;)
HAS_DOCKER := $(shell command -v docker;)
VERSION := ""
DIST := $(CURDIR)/_dist
BUILD := $(CURDIR)/_build
LDFLAGS := "-X main.version=${VERSION}"
BINARY := capctl
CHOCO_SERVER := http://ci.softleader.com.tw:8081/repository/choco/
CHOCO_USER := choco:choco

.PHONY: test
test:
	go test ./... -v

.PHONY: build
build: clean
	mkdir -p $(BUILD)
	go build -o $(BUILD)/$(BINARY) ./cmd/$(BINARY)

.PHONY: dist
dist:
ifndef VERSION
	$(error VERSION is not set)
endif
	mkdir -p $(DIST)
	cp README.md $(DIST) && cp LICENSE $(DIST)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(DIST)/$(BINARY) -a -tags netgo ./cmd/$(BINARY)
	tar -C $(BUILD) -zcvf $(DIST)/$(BINARY)-linux-$(VERSION).tgz $(BINARY) README.md LICENSE
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(DIST)/$(BINARY) -a -tags netgo ./cmd/$(BINARY)
	tar -C $(BUILD) -zcvf $(DIST)/$(BINARY)-darwin-$(VERSION).tgz $(BINARY) README.md LICENSE
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(DIST)/$(BINARY) -a -tags netgo ./cmd/$(BINARY)
	tar -C $(BUILD) -zcvf $(DIST)/$(BINARY)-windows-$(VERSION).tgz $(BINARY) README.md LICENSE

.PHONY: clean
clean:
	rm -rf _*