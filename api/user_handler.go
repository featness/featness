package api

import (
	"encoding/json"
	"fmt"
	"github.com/globoi/featness/api/models"
	"log"
	"net/http"
)

type AutoCompleteUserData struct {
	Name     string
	UserId   string
	ImageUrl string
}

func FindUsersWithIdLike(w http.ResponseWriter, r *http.Request) {
	userid := r.FormValue("name")

	if userid == "" {
		log.Println("Invalid query when auto-completing userid.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users, err := models.FindUsersWithIdLike(userid)

	if err != nil && err.Error() != "not found" {
		log.Println(fmt.Sprintf("Error auto-completing users for term %s (%v).", userid, err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if users == nil {
		w.Write([]byte("[]"))
		return
	}

	userData := []AutoCompleteUserData{}

	for _, user := range *users {
		userData = append(userData, AutoCompleteUserData{user.Name, user.UserId, user.ImageUrl})
	}

	b, err := json.Marshal(userData)

	if err != nil {
		log.Println(fmt.Sprintf("Error converting user to json format (%v).", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(b)
}
