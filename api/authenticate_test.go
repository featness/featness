package api

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tsuru/config"
	"net/http"
	"net/http/httptest"
	"os"
)

func HelperAuth(token string, account string) (string, error) {
	return "my-access-token", nil
}

var _ = Describe("API authenticate Module", func() {

	Context("when StrConfigOrEnv is called", func() {

		It("should have result empty string", func() {
			result := StrConfigOrEnv("key-that-not-exists")
			Expect(result).Should(Equal(""))
		})

		It("should have get var from env", func() {
			os.Setenv("my-test-key", "value-of-my-test-key")
			result := StrConfigOrEnv("my-test-key")
			Expect(result).Should(Equal("value-of-my-test-key"))
		})

		It("should have get uppercase var from env", func() {
			os.Setenv("my-test-key", "")
			os.Setenv("MY-TEST-KEY", "value-of-my-test-key")
			result := StrConfigOrEnv("my-test-key")
			Expect(result).Should(Equal("value-of-my-test-key"))
		})

		It("should have get var from config", func() {
			os.Setenv("MY-TEST-KEY", "")
			config.Set("my-test-key", "value-of-my-test-key")
			result := StrConfigOrEnv("my-test-key")
			Expect(result).Should(Equal("value-of-my-test-key"))
		})

	})

	Context("When authenticate", func() {

		It("should have to generate token", func() {
			config.Set("security_key", "my-secret")
			auth, _ := Authenticate("test-provider", "my-token", "foo@bar.baz", HelperAuth)
			config.Unset("security_key")

			var securityKey string
			token, _ := jwt.Parse(auth, func(t *jwt.Token) ([]byte, error) { return []byte(securityKey), nil })

			Expect(token.Claims["token"]).Should(Equal("my-access-token"))
			Expect(token.Claims["sub"]).Should(Equal("foo@bar.baz"))
			Expect(token.Claims["iss"]).Should(Equal("test-provider"))
			Expect(token.Claims["iat"]).ShouldNot(BeNil())
			Expect(token.Claims["exp"]).ShouldNot(BeNil())
		})

		It("should fail when authenticator fail", func() {
			auth, err := Authenticate("test-provider", "my-token", "foo@bar.baz", func(token string, account string) (string, error) { return "", fmt.Errorf("Any error") })

			Expect(auth).Should(Equal(""))
			Expect(err).Should(MatchError("error authenticating with provider test-provider: Any error"))
		})

		It("should fail when security_key does not in config", func() {
			auth, err := Authenticate("test-provider", "my-token", "foo@bar.baz", HelperAuth)

			Expect(auth).Should(Equal(""))
			Expect(err).Should(MatchError("security key could not be found in configuration file."))
		})

		PIt("should fail when jwt.SignedString throw error", func() {
			// hard to mock to pass
			auth, err := Authenticate("test-provider", "my-token", "foo@bar.baz", HelperAuth)

			Expect(auth).Should(Equal(""))
			Expect(err).Should(MatchError("security key could not be found in configuration file."))
		})

	})

	Context("When AuthenticationRoute", func() {
		It("Should have to set headers", func() {
			loadConfig("../testdata/etc/featness-api1.conf")
			recorder := httptest.NewRecorder()
			request, _ := http.NewRequest("GET", "/authenticate/google", nil)
			request.Header.Add("X-Auth-Data", "foo@bar.com;my-code")
			AuthenticationRoute(recorder, request, "Google", HelperAuth)

			Expect(recorder.Code).Should(Equal(200))
			Expect(recorder.HeaderMap["X-Auth-Token"]).ShouldNot(BeNil())
		})

		It("Should have code 401 when x-auth-data not setted in request", func() {
			recorder := httptest.NewRecorder()
			request, _ := http.NewRequest("GET", "/authenticate/google", nil)
			AuthenticationRoute(recorder, request, "Google", HelperAuth)

			Expect(recorder.Code).Should(Equal(401))
		})

		It("Should have code 401 when x-auth-data not setted in request", func() {
			recorder := httptest.NewRecorder()
			request, _ := http.NewRequest("GET", "/authenticate/google", nil)
			request.Header.Add("X-Auth-Data", "foo@bar.com;my-code")
			AuthenticationRoute(recorder, request, "Google", func(token string, account string) (string, error) { return "", fmt.Errorf("Any error") })

			Expect(recorder.Code).Should(Equal(401))
		})
	})
})
