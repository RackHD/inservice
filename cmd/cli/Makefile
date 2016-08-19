# Build Output Directory
GOOUT ?= $(GOPATH)/bin

# Test Arguments
test_args = -cover -race -trace -randomizeAllSpecs

# variable definitions
SHELL = /bin/bash
PKGPATH = github.com/skunkworxs/inservice-cli
NAME = inservice-cli
COMMITHASH = $(shell git describe --tags --always --dirty)
BUILDDATE = $(shell date -u)
BUILDER = $(shell echo "`git config user.name` <`git config user.email`>")
GOVERSION = $(shell go version)
OSARCH = $(shell uname -sm)
RELEASEVERSION = 0.1

#Flags to pass to cmd/version.go
LDFLAGS = -ldflags "-X '$(PKGPATH)/cmd.binaryName=$(NAME)'\
		    -X '$(PKGPATH)/cmd.buildDate=$(BUILDDATE)'\
		    -X '$(PKGPATH)/cmd.buildUser=$(BUILDER)'\
		    -X '$(PKGPATH)/cmd.commitHash=$(COMMITHASH)'\
		    -X '$(PKGPATH)/cmd.goVersion=$(GOVERSION)'\
		    -X '$(PKGPATH)/cmd.osArch=$(OSARCH)'\
		    -X '$(PKGPATH)/cmd.releaseVersion=$(RELEASEVERSION)'"

# Linter Arguments
#	dupl linter appears to identify errors inaccurately.
lint_args = --vendor --fast --disable=dupl --skip=grpc ./...

default: test

deps:
	go get -u github.com/onsi/ginkgo/ginkgo
	go get -u github.com/onsi/gomega
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install --update
	go get ./...

build: lint
	go build -o $(GOOUT)/$(NAME) $(LDFLAGS) *.go

lint:
	gometalinter $(lint_args)

test: lint
	ginkgo $(test_args)

cover: test
	go tool cover -html=main.coverprofile

watch:
	ginkgo watch $(test_args)
