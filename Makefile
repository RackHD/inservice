# Build Output Directory
GOOUT ?= $(GOPATH)/bin

# Test Arguments
test_args = -cover -race -trace -randomizeAllSpecs

# Linter Arguments
#	dupl linter appears to identify errors inaccurately.
lint_args = --vendor --fast --disable=dupl --skip=grpc ./...

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
	gometalinter --install --update
	go get ./...

build: lint
	go build -o $(GOOUT)/$(NAME) $(LDFLAGS) *.go

run: build
	$(GOOUT)/$(NAME)

lint:
	gometalinter $(lint_args)

test: lint
	ginkgo $(test_args)

cover: test
	go tool cover -html=main.coverprofile

watch:
	ginkgo watch $(test_args)
