package api

import (
	"code.google.com/p/goauth2/oauth"
	"launchpad.net/gocheck"
	"log"
	"os"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { gocheck.TestingT(t) }

type Suite struct {
}

var _ = gocheck.Suite(&Suite{})

func GetGoogleTransport(clientId string, clientSecret string, cacheFilePath string) *oauth.Transport {
	config := &oauth.Config{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  "oob",
		Scope:        "https://www.googleapis.com/auth/plus.profile.emails.read",
		AuthURL:      "https://accounts.google.com/o/oauth2/auth",
		TokenURL:     "https://accounts.google.com/o/oauth2/token",
		TokenCache:   oauth.CacheFile(cacheFilePath),
	}

	code := os.Getenv("GOOGLE_OAUTH_CODE")
	if code == "" {
		url := config.AuthCodeURL("")
		log.Panicf("Visit this URL (%s) to get a code, then put it in an environment variable called GOOGLE_OAUTH_CODE.\n", url)
	}

	transport := &oauth.Transport{Config: config}

	// Try to pull the token from the cache; if this fails, we need to get one.
	token, err := config.TokenCache.Token()
	if err != nil {
		// Exchange the authorization code for an access token.
		// ("Here's the code you gave the user, now give me a token!")
		token, err = transport.Exchange(code)
		if err != nil {
			log.Fatal("Exchange:", err)
		}
	}

	transport.Token = token

	return transport
}
