ORGANIZATION = RackHD
PROJECT = inservice

GOOUT = bin

SHELL = /bin/bash

TTY = $(shell if [ -t 0 ]; then echo "-ti"; fi)

DOCKER_DIR = /go/src/github.com/${ORGANIZATION}/${PROJECT}
DOCKER_IMAGE = rackhd/golang:1.7.0-wheezy
DOCKER_CMD = docker run --rm -v ${PWD}:${DOCKER_DIR} ${TTY} -w ${DOCKER_DIR} ${DOCKER_IMAGE}


# variable definitions
COMMITHASH = $(shell git describe --tags --always --dirty)
BUILDDATE = $(shell date -u)
BUILDER = $(shell echo "`git config user.name` <`git config user.email`>")
GOVERSION = $(shell go version)
OSARCH = $(shell uname -sm)
RELEASEVERSION = 0.1

LINT_ARGS = --vendor --fast --disable=dupl --disable=gotype --disable=gas --disable=gotype --skip=grpc ./...


#Flags to pass to main.go
LDFLAGS = -ldflags "-X 'main.binaryName=${APPLICATION}' \
	  -X 'main.buildDate=${BUILDDATE}' \
	  -X 'main.buildUser=${BUILDER}' \
	  -X 'main.commitHash=${COMMITHASH}' \
	  -X 'main.goVersion=${GOVERSION}' \
	  -X 'main.osArch=${OSARCH}' \
	  -X 'main.releaseVersion=${RELEASEVERSION}' "


#Some tests need to run for 5+ seconds, which trips Ginkgo Slow Test warning
SLOWTEST = 10

.PHONY: shell deps deps-local grpc grpc-local build build-local lint lint-local test test-local release

default: deps grpc test build

coveralls:
	@go get github.com/mattn/goveralls
	@go get github.com/modocache/gover
	@go get golang.org/x/tools/cmd/cover
	@gover
	@goveralls -coverprofile=gover.coverprofile -service=travis-ci

shell:
	@${DOCKER_CMD} /bin/bash

clean:
	@${DOCKER_CMD} make clean-local

clean-local:
	@rm -rf bin vendor

deps:
	@${DOCKER_CMD} make deps-local

deps-local:
	@if ! [ -f glide.lock ]; then glide init --non-interactive; fi
	@glide install

build:
	@${DOCKER_CMD} make build-local

build-local: lint-local
	@go build -o $(GOOUT)/inservice-agent $(LDFLAGS) *.go
	@go build -o $(GOOUT)/inservice-lldp $(LDFLAGS) plugins/lldp/*.go
	@go build -o $(GOOUT)/inservice-catalog-compute $(LDFLAGS) plugins/catalog-compute/*.go
	@go build -o $(GOOUT)/inservice-cli $(LDFLAGS) cmd/cli/*.go


lint:
	@${DOCKER_CMD} make lint-local

lint-local:
	@gometalinter ${LINT_ARGS}

test:
	@${DOCKER_CMD} make test-local

test-local: lint-local
	@ginkgo -r -race -trace -cover -randomizeAllSpecs --slowSpecThreshold=${SLOWTEST}

grpc:
	@${DOCKER_CMD} make grpc-local

grpc-local:
	#Agent's GRPC
	rm -f ./agent/grpc/plugin/plugin.pb.go
	protoc -I ./agent/grpc/plugin ./agent/grpc/plugin/plugin.proto --go_out=plugins=grpc:agent/grpc/plugin
	rm -f ./agent/grpc/host/host.pb.go
	protoc -I ./agent/grpc/host ./agent/grpc/host/host.proto --go_out=plugins=grpc:agent/grpc/host

	#LLDP's GRPC
	rm -f ./plugins/lldp/grpc/lldp/lldp.pb.go
	protoc -I ./plugins/lldp/grpc/lldp ./plugins/lldp/grpc/lldp/lldp.proto --go_out=plugins=grpc:plugins/lldp/grpc/lldp

