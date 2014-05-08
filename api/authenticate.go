package api

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/globoi/featness/api/models"
	"github.com/tsuru/config"
	"labix.org/v2/mgo/bson"
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

	if err != nil {
		return "", fmt.Errorf("Could not get or create user with account %s and token %s (%v)", userAccount, token, err)
	}

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

func IsAuthenticationValid(securityKey []byte) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		header, ok := r.Header["X-Auth-Token"]

		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("X-Auth-Token status was not found.")
			return
		}

		token, err := jwt.Parse(header[0], func(t *jwt.Token) ([]byte, error) { return securityKey, nil })
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(fmt.Sprintf("X-Auth-Token is not a valid token (%v).", err))
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("X-Auth-Token is not a valid token.")
			return
		}

		conn, usersColl, err := models.Users()
		if err != nil {
			log.Println(fmt.Sprintf("Error connecting to the database (%v).", err))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		defer conn.Close()

		sub := token.Claims["sub"].(string)
		user := &models.User{}
		err = usersColl.Find(bson.M{"userid": sub}).One(user)
		if err != nil {
			log.Println(fmt.Sprintf("Could not find user with userId %s (%v).", sub, err))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
