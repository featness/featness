package api

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/globoi/featness/api/models"
	"labix.org/v2/mgo/bson"
	"net/http"
)

func GetUserTeams(w http.ResponseWriter, r *http.Request, token *jwt.Token) {
	sub := token.Claims["sub"].(string)
	user := &models.User{}
	err := models.Users().Find(bson.M{"userid": sub}).One(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	teams, err := models.GetTeamsFor(user.Id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(teams)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func GetAllTeams(w http.ResponseWriter, r *http.Request) {
	teams := &[]models.Team{}
	err := models.Teams().Find(bson.M{}).All(teams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(teams)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(b)
}
