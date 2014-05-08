package models

import (
	"github.com/extemporalgenome/slug"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Team struct {
	Id      bson.ObjectId   `json:"id"        bson:"_id,omitempty"`
	Name    string          `json:"name"`
	Slug    string          `json:"slug"`
	Owner   bson.ObjectId   `json:"owner"`
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
	err = teamsColl.Find(
		bson.M{"$or": []bson.M{
			bson.M{"members": member},
			bson.M{"owner": member},
		}}).All(&teams)

	if err != nil {
		return nil, err
	}

	return teams, nil
}

func FindTeamBySlug(slug string) (*Team, error) {
	conn, teamsColl, err := Teams()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	team := &Team{}
	err = teamsColl.Find(bson.M{"slug": slug}).One(team)

	if err != nil {
		return nil, err
	}

	return team, nil
}

func GetOrCreateTeam(name string, owner *User, members ...*User) (*Team, error) {
	slug := slug.Slug(name)
	team, _ := FindTeamBySlug(slug)

	if team != nil {
		return team, nil
	}

	conn, teamsColl, err := Teams()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	usersBson := make([]bson.ObjectId, len(members))
	for i, v := range members {
		usersBson[i] = v.Id
	}

	err = teamsColl.Insert(
		&Team{bson.NewObjectId(), name, slug, owner.Id, usersBson},
	)

	if err != nil {
		return nil, err
	}

	team, err = FindTeamBySlug(slug)

	if err != nil {
		return nil, err
	}

	return team, nil
}

func IsTeamNameAvailable(name string) (bool, error) {
	conn, teamsColl, err := Teams()
	if err != nil {
		return false, err
	}
	defer conn.Close()

	slug := slug.Slug(name)

	teams := &[]Team{}
	err = teamsColl.Find(bson.M{"slug": slug}).All(teams)

	if err != nil {
		return false, err
	}

	return len(*teams) == 0, nil
}
