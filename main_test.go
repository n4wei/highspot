package main_test

import (
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("End-to-End Integration Tests", func() {
	Context("Running the code with real inputs", func() {
		It("should produce the expected output JSON file", func() {
			highspotCmd := exec.Command("go", "run", "./main.go", "-m", "./test_assets/expected/input.json", "-c", "./test_assets/expected/changes.json", "-o", "./results.json")
			err := highspotCmd.Run()
			Expect(err).ToNot(HaveOccurred())

			diffCmd := exec.Command("diff", "./results.json", "./test_assets/expected/output_compact.json")
			err = diffCmd.Run()
			// a non-zero exit code will result in an error
			// diff exits 0 when there is no difference
			Expect(err).ToNot(HaveOccurred())

			err = os.Remove("./results.json")
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
