# Build Output Directory
GOOUT ?= ./bin

# Test Arguments
test_args = -cover -race -trace -randomizeAllSpecs

# Linter Arguments
#	dupl linter appears to identify errors inaccurately.
lint_args = --vendor --fast --disable=dupl --disable=gotype --disable=gas --skip=grpc ./...

# variable definitions
SHELL = /bin/bash
NAME = inservice-agent
COMMITHASH = $(shell git describe --tags --always --dirty)
BUILDDATE = $(shell date -u)
BUILDER = $(shell echo "`git config user.name` <`git config user.email`>")
GOVERSION = $(shell go version)
OSARCH = $(shell uname -sm)
RELEASEVERSION = 0.1

#Flags to pass to main.go
LDFLAGS = -ldflags "-X 'main.binaryName=$(NAME)'\
		    -X 'main.buildDate=$(BUILDDATE)'\
		    -X 'main.buildUser=$(BUILDER)'\
		    -X 'main.commitHash=$(COMMITHASH)'\
		    -X 'main.goVersion=$(GOVERSION)'\
		    -X 'main.osArch=$(OSARCH)'\
		    -X 'main.releaseVersion=$(RELEASEVERSION)'"


default: test

deps:
	go get -u github.com/onsi/ginkgo/ginkgo
	go get -u github.com/onsi/gomega
	go get -u github.com/alecthomas/gometalinter
	go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
	gometalinter --install
	go get ./...

build: lint
	go build -o $(GOOUT)/inservice-agent $(LDFLAGS) *.go
	go build -o $(GOOUT)/inservice-lldp $(LDFLAGS) plugins/lldp/*.go
	go build -o $(GOOUT)/inservice-catalog-compute $(LDFLAGS) plugins/catalog-compute/*.go
	go build -o $(GOOUT)/inservice-cli $(LDFLAGS) cmd/cli/*.go

run: build
	$(GOOUT)/$(NAME)

lint:
	gometalinter $(lint_args)

test: lint
	ginkgo -r $(test_args)

grpc:
	#Agent's GRPC
	rm -f ./agent/grpc/plugin/plugin.pb.go
	protoc -I ./agent/grpc/plugin ./agent/grpc/plugin/plugin.proto --go_out=plugins=grpc:agent/grpc/plugin
	rm -f ./agent/grpc/host/host.pb.go
	protoc -I ./agent/grpc/host ./agent/grpc/host/host.proto --go_out=plugins=grpc:agent/grpc/host
	
	#LLDP's GRPC
	rm -f ./plugins/lldp/grpc/lldp/lldp.pb.go
	protoc -I ./plugins/lldp/grpc/lldp ./plugins/lldp/grpc/lldp/lldp.proto --go_out=plugins=grpc:plugins/lldp/grpc/lldp

watch:
	ginkgo -r watch $(test_args)
