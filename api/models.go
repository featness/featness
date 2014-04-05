package api

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Team struct {
	Name    string
	Members []string
}

func Teams() *mgo.Collection {
	conn, _ := Conn()
	return conn.C("teams")
}

func GetTeamsFor(member string) ([]Team, error) {
	teams := []Team{}
	err := Teams().Find(bson.M{"members": member}).All(&teams)

	if err != nil {
		return nil, err
	}

	return teams, nil
}
