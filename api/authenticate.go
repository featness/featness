package api

import (
	"code.google.com/p/goauth2/oauth"
	"github.com/dgrijalva/jwt-go"
	"github.com/tsuru/config"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func GetGoogleTransport(clientId string, clientSecret string, cacheFilePath string) *oauth.Transport {
	oauthConfig := &oauth.Config{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
		Scope:        "https://www.googleapis.com/auth/plus.profile.emails.read",
		AuthURL:      "https://accounts.google.com/o/oauth2/auth",
		TokenURL:     "https://accounts.google.com/o/oauth2/token",
		TokenCache:   oauth.CacheFile(cacheFilePath),
	}

	code := os.Getenv("GOOGLE_OAUTH_CODE")
	if code == "" {
		url := oauthConfig.AuthCodeURL("")
		log.Panicf("Visit this URL (%s) to get a code, then put it in an environment variable called GOOGLE_OAUTH_CODE.\n", url)
	}

	transport := &oauth.Transport{Config: oauthConfig}

	// Try to pull the token from the cache; if this fails, we need to get one.
	token, err := oauthConfig.TokenCache.Token()
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

func AuthenticateWithGoogle(w http.ResponseWriter, r *http.Request) {
	authorizationHeader := r.Header.Get("X-Auth-Data")
	if len(authorizationHeader) == 0 {
		// SET STATUS CODE TO 401
		return
	}

	parts := strings.Split(authorizationHeader, ";")
	email, token := parts[0], parts[1]

	clientId, _ := config.GetString("google_client_id")
	clientSecret, _ := config.GetString("google_client_secret")
	cachePath, _ := config.GetString("google_token_cache_path")

	transport := GetGoogleTransport(clientId, clientSecret, cachePath)
	transport.Token.AccessToken = token

	url := "https://www.googleapis.com/oauth2/v1/userinfo"
	clientResponse, err := transport.Client().Get(url)
	defer clientResponse.Body.Close()
	if err != nil {
		// SET STATUS CODE TO 401 and LOG ERROR
		return
	}

	jwtToken := jwt.New(jwt.GetSigningMethod("HS256"))
	jwtToken.Claims["token"] = token
	jwtToken.Claims["sub"] = email
	jwtToken.Claims["iss"] = "Google"
	jwtToken.Claims["iat"] = time.Now().Unix()
	jwtToken.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	securityKey, err := config.GetString("security_key")
	if err != nil {
		// SET STATUS CODE TO 401 and LOG ERROR
		return
	}

	jwtTokenString, _ := jwtToken.SignedString([]byte(securityKey))

	w.Header().Set("X-Auth-Token", jwtTokenString)
	w.WriteHeader(http.StatusOK)
}
