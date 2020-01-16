VERSION=0.1.0

SHELL=sh
PATH:=$(GOPATH)/bin:$(PATH)
BIN_NAME=aws-es-proxy-go
BUILD_DIR=./build
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)
CGO_ENABLED=0

VERSION_PACKAGE=main
GIT_COMMIT?=$(shell git rev-parse --short HEAD)
GIT_DIRTY?=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
GIT_BUILD_DATE=$(shell date --utc '+%Y-%m-%dT%H:%M:%SZ')
GIT_BUILD_GO_VERSION=$(shell go version)
BUILD_LDFLAGS=-s -w \
	-X $(VERSION_PACKAGE).version=$(VERSION) \
	-X $(VERSION_PACKAGE).commit=$(GIT_COMMIT)$(GIT_DIRTY) \
	-X $(VERSION_PACKAGE).buildDate=$(GIT_BUILD_DATE) \
	-X $(VERSION_PACKAGE).buildGoVersion=$(GIT_BUILD_DATE)

# embed version and revision
CURRENT_VERSION=$(shell git log --merges --oneline | perl -ne 'if(m/^.+Merge pull request \#[0-9]+ from .+\/bump-version-([0-9\.]+)/){print $$1;exit}')

.PHONY: run clean install_deps build vet lint fmt test
default: test build

run:
	CGO_ENABLED=$(CGO_ENABLED) go run -installsuffix esproxy -ldflags "$(BUILD_LDFLAGS)" . -config ./example/sample.json

clean:
	rm -rf "$(BUILD_DIR)/*"

install_deps:
	# do nothing

build:
	CGO_ENABLED=$(CGO_ENABLED) go build -o "$(BUILD_DIR)/$(BIN_NAME)" -installsuffix esproxy -ldflags "$(BUILD_LDFLAGS)" .

vet: lint
	go vet ./...

lint:
	# lint: TODO

fmt:
	gofmt -s -l -w $(GOFMT_FILES)

test:
	# test: TODO

#build-docker:
#	docker build .
