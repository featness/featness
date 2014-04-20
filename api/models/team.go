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

func Teams() (*mgo.Session, *mgo.Collection, error) {
	conn, db, err := Conn()
	if err != nil {
		return nil, nil, err
	}
	return conn, db.C("teams"), nil
}

func GetTeamsFor(member bson.ObjectId) ([]Team, error) {
	conn, teamsColl, err := Teams()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	teams := []Team{}
	err = teamsColl.Find(bson.M{"members": member}).All(&teams)

	if err != nil {
		return nil, err
	}

	return teams, nil
}

func GetOrCreateTeam(name string, members ...User) (*Team, error) {
	conn, teamsColl, err := Teams()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	usersBson := make([]bson.ObjectId, len(members))
	for i, v := range members {
		usersBson[i] = v.Id
	}

	_, err = teamsColl.Upsert(
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
