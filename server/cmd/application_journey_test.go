// +build journeyTests

package main_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("The application", func() {
	Context("when pinged", func() {
		var resp *http.Response

		BeforeEach(func() {
			var err error
			resp, err = http.Get("http://app:8080/ping")
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns a HTTP 200 response", func() {
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})
	})
})
