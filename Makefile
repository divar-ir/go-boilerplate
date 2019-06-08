.PHONY: help generate lint fmt dependencies clean check coverage race .remove_empty_dirs

SRCS = $(patsubst ./%,%,$(shell find . -name "*.go" -not -path "*vendor*" -not -path "*.pb.go"))
PACKAGES := $(shell go list ./... | grep -v /vendor)
PROTOS = $(patsubst ./%,%,$(shell find . -name "*.proto"))
PBS = $(patsubst %.proto,%.pb.go,$(patsubst api%,pkg%,$(PROTOS)))
MOCK_PACKAGES = \
	internal/pkg/cache \
	internal/pkg/provider

MOCKED_FILES = $(shell find . -name DOES_NOT_EXIST_FILE $(patsubst %,-or -path "./%/mocks/*.go",$(MOCK_PACKAGES)))
MOCKED_FOLDERS = $(patsubst %,%/mocks,$(MOCK_PACKAGES))

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

generate: $(PBS) $(MOCKED_FILES) $(MOCKED_FOLDERS) | .remove_empty_dirs ## Generate all auto-generated files
.remove_empty_dirs:
	-find . -type d -print | xargs rmdir 2>/dev/null | true

dependencies: | .pre-check-go ## to install the dependencies
	go get -v ./...

clean: ## to remove generated files
	-rm -rf postviewd
	-find . -type d -name mocks -exec rm -rf \{} +

postviewd: $(SRCS) $(PBS) ## Compile postview daemon
	go build -o $@ ./cmd/$@

lint: .bin/golangci-lint ## to lint the files
	.bin/golangci-lint run --config=.golangci-lint.yml ./...

.bin/golangci-lint:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b .bin/ $(LINTER_VERSION)

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

.SECONDEXPANSION:
$(PBS): $$(patsubst %.pb.go,%.proto,$$(patsubst pkg%,api%,$$@)) | .pre-check-go
	protoc -I. --go_out=plugins=grpc:$(GOPATH)/src ./$<

.SECONDEXPANSION:
$(MOCKED_FOLDERS): | .pre-check-go
	cd $(patsubst %/mocks,%,$@) && mockery -all -outpkg mocks -output mocks

.SECONDEXPANSION:
$(MOCKED_FILES): $$(shell find $$(patsubst %/mocks,%,$$(patsubst %/mocks/,%,$$(dir $$@))) -maxdepth 1 -name "*.go") | $(MOCKED_FOLDERS)
	rm -rf $(dir $@)
	cd $(patsubst %/mocks,%,$(patsubst %/mocks/,%,$(dir $@))) && mockery -all -outpkg mocks -output mocks

.pre-check-go:
	go get -v github.com/golang/protobuf/protoc-gen-go
	go get -v github.com/vektra/mockery/.../

# Variables
PROTOC ?= protoc
PROTOC_OPTIONS ?=
LINTER_VERSION = v1.12.5

# Helper Variables

# a variable containing a new line e.g.
# ${\n} would emit a new line
# useful in $(foreach functions
define \n


endef
