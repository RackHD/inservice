package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestInserviceLldp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Inservice LLDP Main Suite")
}
