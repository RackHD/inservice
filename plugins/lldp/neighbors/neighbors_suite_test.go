package neighbors_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestNeighbors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Inservice LLDP Neighbors Suite")
}
