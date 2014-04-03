package api

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/tsuru/config"
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
