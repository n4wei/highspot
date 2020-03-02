package mixtape_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMixtape(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mixtape Suite")
}
