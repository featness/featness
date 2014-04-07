package models

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Team struct {
	Id      bson.ObjectId   `json:"id"        bson:"_id,omitempty"`
	Name    string          `json:"name"`
	Members []bson.ObjectId `json:"members"`
}

func Teams() *mgo.Collection {
	conn, _ := Conn()
	return conn.C("teams")
}

func GetTeamsFor(member bson.ObjectId) ([]Team, error) {
	teams := []Team{}
	err := Teams().Find(bson.M{"members": member}).All(&teams)

	if err != nil {
		return nil, err
	}

	return teams, nil
}

func GetOrCreateTeam(name string, members ...User) (*Team, error) {
	usersBson := make([]bson.ObjectId, len(members))
	for i, v := range members {
		usersBson[i] = v.Id
	}

	teamsColl := Teams()

	_, err := teamsColl.Upsert(
		bson.M{"name": name},
		&Team{bson.NewObjectId(), name, usersBson},
	)

	if err != nil {
		return nil, err
	}

	team := &Team{}
	err = teamsColl.Find(bson.M{"name": name}).One(team)
	return team, nil
}
