package api

import (
	"bytes"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/maraino/go-mock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

type MockClient struct {
	mock.Mock
}

func (c *MockClient) Do(request *http.Request) (int, error) {
	ret := c.Called(request)
	return ret.Int(0), ret.Error(2)
}

var _ = Describe("API facebook authenticate Module", func() {

	Context("when FacebookAuthenticationProvider is called", func() {

		It("should have generate token", func() {
			c := &MockClient{}
			c.When("DoRequest", mock.Any).Return(&http.Response{}, nil).Times(1)
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest("GET", "/authenticate/facebook", nil)
			Expect(err).Should(BeNil())
			request.Header.Add("X-Auth-Data", fmt.Sprintf("heynemann@gmail.com;Bernardo Heynemann;image-url;my-code"))

			AuthenticateWithFacebook(recorder, request)

			Expect(recorder.Code).Should(Equal(200))

			header, ok := recorder.HeaderMap["X-Auth-Token"]
			Expect(ok).Should(Equal(true))
			Expect(header).ShouldNot(BeNil())

			buf := new(bytes.Buffer)
			buf.Write([]byte("my-security-key"))
			key := buf.Bytes()
			token, err := jwt.Parse(header[0], func(t *jwt.Token) ([]byte, error) { return key, nil })

			Expect(token).ShouldNot(BeNil())
			Expect(token.Valid).Should(Equal(true))
			Expect(token.Claims["token"]).ShouldNot(BeNil())
			Expect(token.Claims["sub"]).Should(Equal("heynemann@gmail.com"))
			Expect(token.Claims["iss"]).Should(Equal("Facebook"))
			Expect(token.Claims["iat"]).ShouldNot(BeNil())
			Expect(token.Claims["exp"]).ShouldNot(BeNil())
		})

	})

})
