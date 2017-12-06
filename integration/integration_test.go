package integration_test

import (
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Integration", func() {
	Describe("GET /", func() {
		It("returns 200 OK", func() {
			resp, err := http.Get("http://localhost:8080")
			Expect(err).To(Succeed())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		PIt("shows a list of apps", func() {
			resp, err := http.Get("http://localhost:8080")
			Expect(err).To(Succeed())

			bytes, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(Succeed())
			defer resp.Body.Close()

			str := string(bytes)
			Expect(str).To(ContainSubstring("possum"))
		})
	})
})
