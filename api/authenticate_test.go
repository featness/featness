package api

import (
	"launchpad.net/gocheck"
	"log"
	"os"
)

var _ = gocheck.Suite(&Suite{})

func (s *Suite) TestAuthenticateWithGoogle(c *gocheck.C) {
	clientId := os.Getenv("GOOGLE_CLIENT_ID")
	if clientId == "" {
		log.Panic("Please put your google oauth app client id in an environment variable called GOOGLE_CLIENT_ID.\n")
	}

	secret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if secret == "" {
		log.Panic("Please put your google oauth app client secret in an environment variable called GOOGLE_CLIENT_SECRET.\n")
	}

	transport := GetGoogleTransport(
		clientId,
		secret,
		"/tmp/cache.json",
	)
}
