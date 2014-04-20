package api

import (
	"bytes"
	"code.google.com/p/goauth2/oauth"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tsuru/config"
	"net/http"
	"net/http/httptest"
	"os"
)

func GetGoogleOAuthCode(oauthConfig *oauth.Config) (string, error) {
	code := os.Getenv("GOOGLE_OAUTH_CODE")
	if code == "" {
		url := oauthConfig.AuthCodeURL("")
		return "", fmt.Errorf("Visit this URL (%s) to get a code, then put it in an environment variable called GOOGLE_OAUTH_CODE.\n", url)
	}

	return code, nil
}

var _ = Describe("API google authenticate Module", func() {

	Context("when GoogleAuthenticationProvider is called", func() {

		PIt("should have generate token", func() {
			loadConfig("../testdata/etc/featness-api1.conf")
			clientId := os.Getenv("GOOGLE_CLIENT_ID")
			if clientId == "" {
				fmt.Println("Please put your google oauth app client id in an environment variable called GOOGLE_CLIENT_ID.\n")
			}

			secret := os.Getenv("GOOGLE_CLIENT_SECRET")
			if secret == "" {
				fmt.Println("Please put your google oauth app client secret in an environment variable called GOOGLE_CLIENT_SECRET.\n")
			}

			config.Set("google_client_id", clientId)
			config.Set("google_client_secret", secret)
			config.Set("google_token_cache_path", "/tmp/cache.json")

			oauthConfig, err := GetGoogleOAuthConfig()
			if err != nil {
				fmt.Println(err.Error())
			}

			code, err := GetGoogleOAuthCode(oauthConfig)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			recorder := httptest.NewRecorder()

			request, err := http.NewRequest("GET", "/authenticate/google", nil)
			Expect(err).Should(BeNil())
			request.Header.Add("X-Auth-Data", fmt.Sprintf("heynemann@gmail.com;%s", code))

			AuthenticateWithGoogle(recorder, request)

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
			Expect(token.Claims["iss"]).Should(Equal("Google"))
			Expect(token.Claims["iat"]).ShouldNot(BeNil())
			Expect(token.Claims["exp"]).ShouldNot(BeNil())
		})

	})

})
