package api

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Healthcheck Handler", func() {
	It("should return WORKING as result", func() {
		recorder := httptest.NewRecorder()
		request, err := http.NewRequest("GET", "/healthcheck", nil)
		Expect(err).ShouldNot(HaveOccurred())

		Healthcheck(recorder, request)
		Expect(recorder.Code).Should(Equal(http.StatusOK))

		body, err := ioutil.ReadAll(recorder.Body)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(string(body)).Should(Equal("WORKING"))
	})
})
