package api

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/globoi/featness/api/models"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"strings"
)

func GetUserTeams(w http.ResponseWriter, r *http.Request, token *jwt.Token) {
	conn, usersColl, err := models.Users()
	if err != nil {
		log.Println(fmt.Sprintf("Error connecting to the database (%v).", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	sub := token.Claims["sub"].(string)
	user := &models.User{}
	err = usersColl.Find(bson.M{"userid": sub}).One(user)
	if err != nil {
		log.Println(fmt.Sprintf("Could not find user with userId %s (%v).", sub, err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	teams, err := models.GetTeamsFor(user.Id)

	if err != nil {
		log.Println(fmt.Sprintf("Error retrieving user %s's teams (%v).", sub, err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(teams)

	if err != nil {
		log.Println(fmt.Sprintf("Error converting user %s's teams to json format (%v).", sub, err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func GetAllTeams(w http.ResponseWriter, r *http.Request) {
	conn, teamsColl, err := models.Teams()
	if err != nil {
		log.Println(fmt.Sprintf("Error connecting to the database (%v).", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	teams := &[]models.Team{}
	err = teamsColl.Find(bson.M{}).All(teams)
	if err != nil {
		log.Println(fmt.Sprintf("Error retrieving all teams (%v).", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(teams)

	if err != nil {
		log.Println(fmt.Sprintf("Error converting all teams to json format (%v).", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func IsTeamNameAvailable(w http.ResponseWriter, r *http.Request, token *jwt.Token) {
	name := r.FormValue("name")
	if name == "" {
		log.Println("Invalid team name when finding if team name available.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	isTeamNameAvailable, err := models.IsTeamNameAvailable(name)
	if err != nil {
		log.Println(fmt.Sprintf("Error finding all teams with name that matches '%s' (%v).", name, err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(isTeamNameAvailable)

	if err != nil {
		log.Println(fmt.Sprintf("Error converting result to json format (%v).", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func LoadTeam(w http.ResponseWriter, r *http.Request, token *jwt.Token) {
	teamId := r.URL.Query().Get(":teamId")

	team, err := models.FindTeamBySlug(teamId)
	if err != nil {
		log.Println(fmt.Sprintf("Error loading team with slug '%s' (%v).", teamId, err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	owner, err := team.Owner.MarshalJSON()
	if err != nil {
		log.Println(fmt.Sprintf("Error serializing team owner with slug '%s' (%v).", teamId, err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	members := []string{}

	for _, member := range team.Members {
		memberId, err := member.MarshalJSON()
		if err != nil {
			log.Println(fmt.Sprintf("Error serializing team user with slug '%s' (%v).", teamId, err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		members = append(members, string(memberId))
	}

	b, err := json.Marshal(
		map[string]interface{}{
			"slug":    team.Slug,
			"name":    team.Name,
			"members": members,
			"owner":   owner,
		},
	)

	if err != nil {
		log.Println(fmt.Sprintf("Error converting result to json format (%v).", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func CreateTeam(w http.ResponseWriter, r *http.Request, token *jwt.Token) {
	name := r.FormValue("name")
	owner := r.FormValue("owner")
	users := r.FormValue("users")

	loadedOwner, err := models.GetUserByUserId(owner)
	if err != nil {
		log.Println(fmt.Sprintf("Invalid owner when saving new team (%v).", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	loadedUsers := []*models.User{}
	if users != "" {
		userIds := strings.Split(users, ",")
		for _, userId := range userIds {
			user, err := models.GetUserByUserId(userId)
			if err != nil {
				log.Println(fmt.Sprintf("Invalid user (%s) when saving new team (%v).", userId, err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			loadedUsers = append(loadedUsers, user)
		}
	}

	team, err := models.GetOrCreateTeam(name, loadedOwner, loadedUsers...)

	if err != nil {
		log.Println(fmt.Sprintf("Error creating new team (%v).", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write([]byte(team.Slug))
}
