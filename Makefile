.PHONY: help generate lint fmt dependencies clean check coverage race .remove_empty_dirs .pre-check-go

SRCS = $(patsubst ./%,%,$(shell find . -name "*.go" -not -path "*vendor*" -not -path "*.pb.go"))
PACKAGES := $(shell go list ./... | grep -v /vendor)
PROTOS = $(patsubst ./%,%,$(shell find . -name "*.proto"))
PBS = $(patsubst %.proto,%.pb.go,$(patsubst api%,pkg%,$(PROTOS)))
MOCK_PACKAGES = \
	internal/pkg/provider \
	internal/pkg/metrics

MOCKED_FILES = $(shell find . -name DOES_NOT_EXIST_FILE $(patsubst %,-or -path "./%/mocks/*.go",$(MOCK_PACKAGES)))
MOCKED_FOLDERS = $(patsubst %,%/mocks,$(MOCK_PACKAGES))

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

generate: $(PBS) $(MOCKED_FILES) $(MOCKED_FOLDERS) cmd/postviewd/wire_gen.go | .remove_empty_dirs ## Generate all auto-generated files
.remove_empty_dirs:
	-find . -type d -print | xargs rmdir 2>/dev/null | true

dependencies: | .pre-check-go .bin/golangci-lint ## to install the dependencies
	go mod download

clean: ## to remove generated files
	-rm -rf postviewd
	-find . -type d -name mocks -exec rm -rf \{} +

postviewd: $(SRCS) $(PBS) ## Compile postview daemon
	go build -o $@ -ldflags="$(LD_FLAGS)" ./cmd/$@

docker: ## to build docker image
	$(DOCKER) build -t $(IMAGE_NAME):$(IMAGE_VERSION) .

push: docker ## to push docker image to registry
	$(DOCKER) push $(IMAGE_NAME):$(VERSION)

push-production: ## to tag and push :production tag on docker image
	$(DOCKER) pull $(IMAGE_NAME):$(IMAGE_VERSION)
	$(DOCKER) tag $(IMAGE_NAME):$(IMAGE_VERSION) $(IMAGE_NAME):production
	$(DOCKER) push $(IMAGE_NAME):production

deploy: ## to deploy it on kubernetes
	kubectl --namespace divar-review patch deployment/postview -p='{"spec":{"template":{"spec":{"containers":[{"name":"postview","imagePullPolicy":"IfNotPresent"}]}}}}' || echo "No Need To Patch Config"
	kubectl --namespace divar-review set image deployment/postview postview=$(IMAGE_NAME):$(VERSION)

lint: .bin/golangci-lint ## to lint the files
	.bin/golangci-lint run --config=.golangci-lint.yml ./...

.bin/golangci-lint:
	if [ -z "$$(which golangci-lint)" ]; then curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b .bin/ $(LINTER_VERSION); else mkdir -p .bin; ln -s "$$(which golangci-lint)" $@; fi

fmt: ## to run `go fmt` on all source code
	gofmt -s -w $(SRCS)

check: | generate ## Run tests
	go test ./...

race: | generate ## to run data race detector
	go test -timeout 30s -race ./...

coverage: coverage.cover coverage.html ## to run tests and generate test coverage data
	gocov convert $< | gocov report

coverage.html: coverage.cover
	go tool cover -html=$< -o $@

coverage.cover: $(SRCS) $(PBS) Makefile | generate
	-rm -rfv .coverage
	mkdir -p .coverage
	$(foreach pkg,$(PACKAGES),go test -timeout 30s -short -covermode=count -coverprofile=.coverage/$(subst /,-,$(pkg)).cover $(pkg)${\n})
	echo "mode: count" > $@
	grep -h -v "^mode:" .coverage/*.cover >> $@

cmd/postviewd/wire_gen.go: cmd/postviewd/container.go
	wire ./cmd/postviewd

.SECONDEXPANSION:
$(PBS): $$(patsubst %.pb.go,%.proto,$$(patsubst pkg%,api%,$$@)) | .pre-check-go
	$(PROTOC) $(PROTOC_OPTIONS) --go_out=plugins=grpc:$(GOPATH)/src ./$<

.SECONDEXPANSION:
$(MOCKED_FOLDERS): | .pre-check-go
	cd $(patsubst %/mocks,%,$@) && mockery -all -outpkg mocks -output mocks

.SECONDEXPANSION:
$(MOCKED_FILES): $$(shell find $$(patsubst %/mocks,%,$$(patsubst %/mocks/,%,$$(dir $$@))) -maxdepth 1 -name "*.go") | $(MOCKED_FOLDERS)
	rm -rf $(dir $@)
	cd $(patsubst %/mocks,%,$(patsubst %/mocks/,%,$(dir $@))) && mockery -all -outpkg mocks -output mocks

.pre-check-go: 
	if [ -z "$$(which protoc-gen-go)" ]; then go get -v github.com/golang/protobuf/protoc-gen-go; fi
	if [ -z "$$(which mockery)" ]; then go get -v github.com/vektra/mockery/cmd/mockery; fi
	if [ -z "$$(which gocov)" ]; then go get -v github.com/axw/gocov/gocov; fi
	if [ -z "$$(which wire)" ]; then go get -v github.com/google/wire/cmd/wire; fi

# Variables
ROOT := github.com/cafebazaar/go-boilerplate

PROTOC ?= protoc
PROTOC_OPTIONS ?= -I.
LINTER_VERSION = v1.12.5
GIT ?= git
DOCKER ?= docker
COMMIT := $(shell $(GIT) rev-parse HEAD)
CI_COMMIT_TAG ?=
VERSION ?= $(strip $(if $(CI_COMMIT_TAG),$(CI_COMMIT_TAG),$(shell $(GIT) describe --tag 2> /dev/null || echo "$(COMMIT)")))
BUILD_TIME := $(shell LANG=en_US date +"%F_%T_%z")
LD_FLAGS := -X $(ROOT)/pkg/postview.Version=$(VERSION) -X $(ROOT)/pkg/postview.Commit=$(COMMIT) -X $(ROOT)/pkg/postview.BuildTime=$(BUILD_TIME)
IMAGE_NAME ?= registry.cafebazaar.ir:5000/arcana261/golang-boilerplate
IMAGE_VERSION ?= $(VERSION)

# Helper Variables

# a variable containing a new line e.g.
# ${\n} would emit a new line
# useful in $(foreach functions
define \n


endef
