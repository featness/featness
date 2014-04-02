package api

import (
	"bytes"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"launchpad.net/gocheck"
	"net/http"
	"net/http/httptest"
	"os"
)

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

	transport := GetGoogleTransport(
		clientId,
		secret,
		"/tmp/cache.json",
	)

	recorder := httptest.NewRecorder()

	// X-AUTH-DATA="heynemann@gmail.com;qwi9129349124912"
	request, err := http.NewRequest("GET", "/authenticate/google", nil)
	c.Assert(err, gocheck.IsNil)
	request.Header.Add("X-Auth-Data", fmt.Sprintf("heynemann@gmail.com;%s", transport.Token.AccessToken))

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
	c.Assert(token.Claims["token"], gocheck.Equals, transport.Token.AccessToken)
	c.Assert(token.Claims["sub"], gocheck.Equals, "heynemann@gmail.com")
	c.Assert(token.Claims["iss"], gocheck.Equals, "Google")
	c.Assert(token.Claims["iat"], gocheck.NotNil)
	c.Assert(token.Claims["exp"], gocheck.NotNil)
}
