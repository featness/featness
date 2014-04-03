package api

import (
	"bytes"
	"code.google.com/p/goauth2/oauth"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/tsuru/config"
	"launchpad.net/gocheck"
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

func (s *Suite) TestAuthenticateWithGoogle(c *gocheck.C) {
	loadConfig("../testdata/etc/featness-api1.conf")
	clientId := os.Getenv("GOOGLE_CLIENT_ID")
	if clientId == "" {
		c.Fatal("Please put your google oauth app client id in an environment variable called GOOGLE_CLIENT_ID.\n")
	}

	secret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if secret == "" {
		c.Fatal("Please put your google oauth app client secret in an environment variable called GOOGLE_CLIENT_SECRET.\n")
	}

	config.Set("google_client_id", clientId)
	config.Set("google_client_secret", secret)
	config.Set("google_token_cache_path", "/tmp/cache.json")

	oauthConfig, err := GetGoogleOAuthConfig()
	if err != nil {
		c.Fatal(err)
	}

	code, err := GetGoogleOAuthCode(oauthConfig)
	if err != nil {
		c.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	request, err := http.NewRequest("GET", "/authenticate/google", nil)
	c.Assert(err, gocheck.IsNil)
	request.Header.Add("X-Auth-Data", fmt.Sprintf("heynemann@gmail.com;%s", code))

	AuthenticateWithGoogle(recorder, request)

	c.Assert(recorder.Code, gocheck.Equals, http.StatusOK)

	header, ok := recorder.HeaderMap["X-Auth-Token"]
	c.Assert(ok, gocheck.Equals, true)
	c.Assert(header, gocheck.NotNil)

	buf := new(bytes.Buffer)
	buf.Write([]byte("my-security-key"))
	key := buf.Bytes()
	token, err := jwt.Parse(header[0], func(t *jwt.Token) ([]byte, error) { return key, nil })

	c.Assert(token, gocheck.NotNil)

	c.Assert(token.Valid, gocheck.Equals, true)
	c.Assert(token.Claims["token"], gocheck.NotNil)
	c.Assert(token.Claims["sub"], gocheck.Equals, "heynemann@gmail.com")
	c.Assert(token.Claims["iss"], gocheck.Equals, "Google")
	c.Assert(token.Claims["iat"], gocheck.NotNil)
	c.Assert(token.Claims["exp"], gocheck.NotNil)
}
