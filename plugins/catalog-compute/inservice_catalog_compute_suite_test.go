package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestInserviceCatalogCompute(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "InserviceCatalogCompute Suite")
}
