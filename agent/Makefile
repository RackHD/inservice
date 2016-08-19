
# Test Arguments
test_args = -cover -race -trace -randomizeAllSpecs

# Linter Arguments
#	dupl linter appears to identify errors inaccurately.
lint_args = --vendor --fast --disable=dupl --disable=gas --skip=grpc ./...

default: test

deps:
	go get -u github.com/onsi/ginkgo/ginkgo
	go get -u github.com/onsi/gomega
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install --update
	go get ./...

build: lint
	go build *.go

lint:
	gometalinter $(lint_args)

test: lint
	ginkgo $(test_args)

cover: test
	go tool cover -html=plugins.coverprofile

watch:
	ginkgo watch $(test_args)

grpc:
	rm -f ./grpc/plugin/plugin.pb.go
	protoc -I ./grpc/plugin ./grpc/plugin/plugin.proto --go_out=plugins=grpc:grpc/plugin
	rm -f ./grpc/host/host.pb.go
	protoc -I ./grpc/host ./grpc/host/host.proto --go_out=plugins=grpc:grpc/host

.PHONY: grpc
