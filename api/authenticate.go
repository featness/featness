package api

import (
	"code.google.com/p/goauth2/oauth"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/tsuru/config"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func GetGoogleTransport(clientId string, clientSecret string, cacheFilePath string) (*oauth.Transport, error) {
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
		msg := fmt.Sprintf("Visit this URL (%s) to get a code, then put it in an environment variable called GOOGLE_OAUTH_CODE.\n", url)
		log.Println(msg)
		return nil, fmt.Errorf(msg)
	}

	transport := &oauth.Transport{Config: oauthConfig}

	// Try to pull the token from the cache; if this fails, we need to get one.
	token, err := oauthConfig.TokenCache.Token()
	if err != nil {
		// Exchange the authorization code for an access token.
		token, err = transport.Exchange(code)
		if err != nil {
			msg := fmt.Sprintf("Error when trying to get a new token with Google (%v).\n", err)
			log.Println(msg)
			return nil, fmt.Errorf(msg)
		}
	}

	transport.Token = token

	return transport, nil
}

func StrConfigOrEnv(key string) string {
	value := os.Getenv(key)

	if value == "" {
		value = os.Getenv(strings.ToUpper(key))
	}
	if value == "" {
		value, _ = config.GetString(key)
	}

	return value
}

func AuthenticateWithGoogle(w http.ResponseWriter, r *http.Request) {
	authorizationHeader := r.Header.Get("X-Auth-Data")
	if len(authorizationHeader) == 0 {
		log.Println("Authorization header was not found in request.")
		// SET STATUS CODE TO 401
		return
	}

	parts := strings.Split(authorizationHeader, ";")
	email, token := parts[0], parts[1]

	clientId := StrConfigOrEnv("google_client_id")
	clientSecret := StrConfigOrEnv("google_client_secret")
	cachePath := StrConfigOrEnv("google_token_cache_path")

	if clientId == "" || clientSecret == "" || cachePath == "" {
		log.Println("Client ID, or Secret or Cache Path were not found in config (or environment variable).")
		// SET STATUS CODE TO 401 and LOG ERROR
		return
	}

	transport, err := GetGoogleTransport(clientId, clientSecret, cachePath)
	if err != nil {
		log.Println(fmt.Sprintf("Google transport could not be configured: %v", err))
		// SET STATUS CODE TO 401 and LOG ERROR
		return
	}
	transport.Token.AccessToken = token

	url := "https://www.googleapis.com/oauth2/v1/userinfo"
	clientResponse, err := transport.Client().Get(url)
	defer clientResponse.Body.Close()
	if err != nil {
		log.Println(fmt.Sprintf("Access token was invalid: %v.", err))
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
		log.Println("Security key could not be found in configuration file.")
		// SET STATUS CODE TO 401 and LOG ERROR
		return
	}

	jwtTokenString, _ := jwtToken.SignedString([]byte(securityKey))

	w.Header().Set("X-Auth-Token", jwtTokenString)
	w.WriteHeader(http.StatusOK)
}
