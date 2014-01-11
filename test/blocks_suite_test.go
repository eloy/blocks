package blocks_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestBlocks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Blocks Suite")
}
