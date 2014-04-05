package api

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

func GetAllTeams(w http.ResponseWriter, r *http.Request, token *jwt.Token) {
	sub := token.Claims["sub"].(string)

	teams, err := GetTeamsFor(sub)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(teams)
	fmt.Println(err, b, teams)
	w.Write(b)
}
