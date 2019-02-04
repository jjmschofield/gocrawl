package pages_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPages(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pages Suite")
}
