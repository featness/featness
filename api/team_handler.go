package api

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

func GetAllTeams(w http.ResponseWriter, r *http.Request, token *jwt.Token) {
	sub := token.Claims["sub"].(string)
	w.Write([]byte(sub))
}
