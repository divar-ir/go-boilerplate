.PHONY: help generate

SRCS = $(patsubst ./%,%,$(shell find . -name "*.go" -not -path "*vendor*" -not -path "*.pb.go"))
PROTOS = $(patsubst ./%,%,$(shell find . -name "*.proto"))
PBS = $(patsubst %.proto,%.pb.go,$(patsubst api%,pkg%,$(PROTOS)))

$(info $(PBS))

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

generate: $(PBS) ## Generate all auto-generated files

postviewd: $(SRCS) $(PBS) ## Compile postview daemon
	go build -o $@ ./cmd/$@

.SECONDEXPANSION:
$(PBS): $$(patsubst %.pb.go,%.proto,$$(patsubst pkg%,api%,$$@)) | .pre-check-go
	protoc -I. --go_out=plugins=grpc:$(GOPATH)/src ./$<

.pre-check-go:
	go get -v github.com/golang/protobuf/protoc-gen-go

# Variables
PROTOC ?= protoc
PROTOC_OPTIONS ?=
