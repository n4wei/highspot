package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHighspot(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Highspot Suite")
}
