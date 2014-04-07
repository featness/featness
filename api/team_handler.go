package api

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"labix.org/v2/mgo/bson"
	"net/http"
)

func GetAllTeams(w http.ResponseWriter, r *http.Request, token *jwt.Token) {
	sub := token.Claims["sub"].(string)
	user := &User{}
	err := Users().Find(bson.M{"userID": sub}).One(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	teams, err := GetTeamsFor(user.Id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(teams)
	fmt.Println(err, b, teams)
	w.Write(b)
}
