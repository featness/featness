package api

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
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

func AuthenticationRoute(w http.ResponseWriter, r *http.Request, providerName string, authenticator AuthenticationProvider) {
	authorizationHeader := r.Header.Get("X-Auth-Data")
	if len(authorizationHeader) == 0 {
		log.Println("authorization header was not found in request.")
		// SET STATUS CODE TO 401
		return
	}

	parts := strings.Split(authorizationHeader, ";")
	userAccount, token := parts[0], parts[1]

	token, err := Authenticate(providerName, token, userAccount, authenticator)

	if err != nil {
		log.Println(err)
		// SET STATUS CODE TO 401
		return
	}

	w.Header().Set("X-Auth-Token", token)
	w.WriteHeader(http.StatusOK)
}
