// +build unitTests

package api_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/batect/abacus/server/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ping endpoint", func() {
	Context("when invoked", func() {
		var resp *httptest.ResponseRecorder

		BeforeEach(func() {
			resp = httptest.NewRecorder()
			api.Ping(resp, nil)
		})

		It("returns a HTTP 200 response", func() {
			Expect(resp.Code).To(Equal(http.StatusOK))
		})

		It("returns 'pong' in the response body", func() {
			Expect(resp.Body.String()).To(Equal("pong"))
		})
	})
})
