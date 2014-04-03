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

func GetGoogleOAuthConfig() (*oauth.Config, error) {
	clientId := StrConfigOrEnv("google_client_id")
	clientSecret := StrConfigOrEnv("google_client_secret")

	if clientId == "" || clientSecret == "" {
		return nil, fmt.Errorf("Client ID or Client Secret were not found in config (or environment variable).")
	}

	oauthConfig := &oauth.Config{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
		Scope:        "https://www.googleapis.com/auth/plus.profile.emails.read",
		AuthURL:      "https://accounts.google.com/o/oauth2/auth",
		TokenURL:     "https://accounts.google.com/o/oauth2/token",
		TokenCache:   nil,
	}

	return oauthConfig, nil
}

func GetGoogleTransport(token string) (*oauth.Transport, error) {
	oauthConfig, err := GetGoogleOAuthConfig()
	if err != nil {
		return nil, err
	}

	transport := &oauth.Transport{Config: oauthConfig}

	transport.Token = &oauth.Token{}
	transport.Token.AccessToken = token

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

type AuthenticationProvider func(token string, account string) (string, error)

func GoogleAuthenticationProvider(token string, account string) (string, error) {
	transport, err := GetGoogleTransport(token)
	if err != nil {
		return "", fmt.Errorf("Google transport could not be configured: %v", err)
	}

	url := "https://www.googleapis.com/oauth2/v1/userinfo"
	clientResponse, err := transport.Client().Get(url)
	defer clientResponse.Body.Close()

	if err != nil || clientResponse.Status == "401" {
		return "", fmt.Errorf("access token was invalid: %v.", err)
	}

	// TODO: Verify if token e-mail account is the same as the e-mail passed

	return transport.Token.AccessToken, nil
}

func Authenticate(provider string, token string, email string, authenticator AuthenticationProvider) (string, error) {
	token, err := authenticator(token, email)

	if err != nil {
		return "", fmt.Errorf("error authenticating with provider %s: %v", provider, err)
	}

	jwtToken := jwt.New(jwt.GetSigningMethod("HS256"))
	jwtToken.Claims["token"] = token
	jwtToken.Claims["sub"] = email
	jwtToken.Claims["iss"] = provider
	jwtToken.Claims["iat"] = time.Now().Unix()
	jwtToken.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	securityKey, err := config.GetString("security_key")
	if err != nil {
		return "", fmt.Errorf("security key could not be found in configuration file.")
	}

	jwtTokenString, _ := jwtToken.SignedString([]byte(securityKey))

	if err != nil {
		return "", fmt.Errorf("security Token could not be generated (%v).", err)
	}

	return jwtTokenString, nil
}

func AuthenticateWithGoogle(w http.ResponseWriter, r *http.Request) {
	authorizationHeader := r.Header.Get("X-Auth-Data")
	if len(authorizationHeader) == 0 {
		log.Println("authorization header was not found in request.")
		// SET STATUS CODE TO 401
		return
	}

	parts := strings.Split(authorizationHeader, ";")
	email, token := parts[0], parts[1]

	token, err := Authenticate("Google", token, email, GoogleAuthenticationProvider)

	if err != nil {
		log.Println(err)
		// SET STATUS CODE TO 401
		return
	}

	w.Header().Set("X-Auth-Token", token)
	w.WriteHeader(http.StatusOK)
}
