package api

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/globoi/featness/api/models"
	"github.com/tsuru/config"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

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

func Authenticate(provider string, token string, userAccount string, name string, imageUrl string, authenticator AuthenticationProvider) (string, error) {
	token, err := authenticator(token, userAccount)

	if err != nil {
		return "", fmt.Errorf("error authenticating with provider %s: %v", provider, err)
	}

	user, err := models.GetOrCreateUser(provider, token, userAccount, name, imageUrl)

	jwtToken := jwt.New(jwt.GetSigningMethod("HS256"))
	jwtToken.Claims["token"] = user.AccessToken
	jwtToken.Claims["sub"] = user.UserId
	jwtToken.Claims["iss"] = user.Provider
	jwtToken.Claims["iat"] = time.Now().Unix()
	jwtToken.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	securityKey, err := config.GetString("security_key")
	if err != nil {
		return "", fmt.Errorf("security key could not be found in configuration file.")
	}

	jwtTokenString, err := jwtToken.SignedString([]byte(securityKey))
	if err != nil {
		return "", fmt.Errorf("security Token could not be generated (%v).", err)
	}

	return jwtTokenString, nil
}

func AuthenticationRoute(w http.ResponseWriter, r *http.Request, providerName string, authenticator AuthenticationProvider) {
	authorizationHeader := r.Header.Get("X-Auth-Data")
	if len(authorizationHeader) == 0 {
		log.Println("authorization header was not found in request.")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	parts := strings.Split(authorizationHeader, ";")
	userAccount, name, imageUrl, token := parts[0], parts[1], parts[2], parts[3]

	token, err := Authenticate(providerName, token, userAccount, name, imageUrl, authenticator)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("X-Auth-Token", token)
	w.WriteHeader(http.StatusOK)
}
