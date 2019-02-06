package counters_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCounters(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Counters Suite")
}
